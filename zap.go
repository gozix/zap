package zap

import (
	"fmt"
	"os"
	"strings"

	glueBundle "github.com/gozix/glue/v2"
	viperBundle "github.com/gozix/viper/v2"

	"github.com/sarulabs/di/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// BundleName is default definition name.
	BundleName = "zap"

	// ArgCoreType is argument name.
	ArgCoreType = "zap.core.type"

	// TagCoreFactory is factory tag.
	TagCoreFactory = "zap.core.factory"

	// DefStreamCoreFactory is definition name.
	DefStreamCoreFactory = "zap.core.stream"
)

type (
	// Bundle implements the glue.Bundle interface.
	Bundle struct{}

	// CoreFactory is type alias of core.Factory.
	CoreFactory = func(path string) (zapcore.Core, error)

	// fieldsConf is logger conf fieldsConf struct.
	fieldsConf []struct {
		Key, Value string
	}
)

// Type checking.
var _ glueBundle.Bundle = (*Bundle)(nil)

// NewBundle create bundle instance.
func NewBundle() *Bundle {
	return new(Bundle)
}

// Key implements the glue.Bundle interface.
func (b *Bundle) Name() string {
	return BundleName
}

// Build implements the glue.Bundle interface.
func (b *Bundle) Build(builder *di.Builder) (err error) {
	var defs = make([]di.Def, 0, 2)
	if !builder.IsDefined(BundleName) {
		defs = append(defs, b.defBundle())
	}

	if !builder.IsDefined(DefStreamCoreFactory) {
		defs = append(defs, b.defStreamCoreFactory())
	}

	return builder.Add(defs...)
}

// DependsOn implements the glue.DependsOn interface.
func (b *Bundle) DependsOn() []string {
	return []string{viperBundle.BundleName}
}

// defBundle is definition getter.
func (b *Bundle) defBundle() di.Def {
	return di.Def{
		Name: BundleName,
		Build: func(ctn di.Container) (_ interface{}, err error) {
			var cfg *viper.Viper
			if err = ctn.Fill(viperBundle.BundleName, &cfg); err != nil {
				return nil, err
			}

			var factories map[string]CoreFactory
			if factories, err = b.factories(ctn); err != nil {
				return nil, err
			}

			var cores []zapcore.Core
			if cores, err = b.cores(cfg, factories); err != nil {
				return nil, err
			}

			var options []zap.Option
			if options, err = b.options(cfg); err != nil {
				return nil, err
			}

			var logger = zap.New(zapcore.NewTee(cores...), options...)
			zap.ReplaceGlobals(logger)
			zap.RedirectStdLog(logger)

			return logger, nil
		},
		Close: func(obj interface{}) (err error) {
			return obj.(*zap.Logger).Sync()
		},
	}
}

// defStreamCoreFactory is definition getter.
func (b *Bundle) defStreamCoreFactory() di.Def {
	return di.Def{
		Name: DefStreamCoreFactory,
		Tags: []di.Tag{{
			Name: TagCoreFactory,
			Args: map[string]string{
				ArgCoreType: "stream",
			},
		}},
		Build: func(ctn di.Container) (interface{}, error) {
			return func(path string) (_ zapcore.Core, err error) {
				var cfg *viper.Viper
				if err = ctn.Fill(viperBundle.BundleName, &cfg); err != nil {
					return nil, err
				}

				var (
					key   = strings.Split(path, ".")[0] + ".development"
					eConf = zap.NewProductionEncoderConfig()
				)

				if cfg.IsSet(key) && cfg.GetBool(key) {
					eConf = zap.NewDevelopmentEncoderConfig()
				}

				key = path + ".time_encoder"
				if cfg.IsSet(key) {
					if err = eConf.EncodeTime.UnmarshalText([]byte(cfg.GetString(key))); err != nil {
						return nil, err
					}
				}

				key = path + ".message_key"
				if cfg.IsSet(key) {
					eConf.MessageKey = cfg.GetString(key)
				}

				var encoding = "json"
				key = path + ".encoding"
				if cfg.IsSet(key) {
					encoding = cfg.GetString(key)
				}

				var encoder zapcore.Encoder
				switch encoding {
				case "json":
					encoder = zapcore.NewJSONEncoder(eConf)
				case "console":
					encoder = zapcore.NewConsoleEncoder(eConf)
				default:
					return nil, fmt.Errorf(`encoding "%s" is not supported`, encoding)
				}

				var level = zap.NewAtomicLevel()
				key = path + ".level"
				if cfg.IsSet(key) {
					if err = level.UnmarshalText([]byte(cfg.GetString(key))); err != nil {
						return nil, err
					}
				}

				return zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level), nil
			}, nil
		},
	}
}

// factories is factories getter.
func (b *Bundle) factories(ctn di.Container) (_ map[string]CoreFactory, err error) {
	var factories = make(map[string]CoreFactory, 4)
	for name, def := range ctn.Definitions() {
		for _, tag := range def.Tags {
			if tag.Name != TagCoreFactory {
				continue
			}

			var coreType, ok = tag.Args[ArgCoreType]
			if !ok {
				return nil, fmt.Errorf(
					`core definition "%s" don't have required argument "%s"`, def.Name, ArgCoreType,
				)
			}

			if _, ok := factories[coreType]; ok {
				return nil, fmt.Errorf(`core type "%s" factory already defined`, coreType)
			}

			var factory CoreFactory
			if err = ctn.Fill(name, &factory); err != nil {
				return nil, err
			}

			factories[coreType] = factory
			break
		}
	}

	return factories, nil
}

// cores is cores getter.
func (b *Bundle) cores(cfg *viper.Viper, factories map[string]CoreFactory) (_ []zapcore.Core, err error) {
	var (
		defs  = cfg.GetStringMap(BundleName + ".cores")
		cores = make([]zapcore.Core, 0, 4)
	)

	for def := range defs {
		var path = fmt.Sprintf("%s.cores.%s", BundleName, def)
		if cfg.IsSet(path + ".type") {
			var coreType = cfg.GetString(path + ".type")
			if _, ok := factories[coreType]; !ok {
				return nil, fmt.Errorf(`core factory with type "%s" is not defined`, coreType)
			}

			var c zapcore.Core
			if c, err = factories[coreType](path); err != nil {
				return nil, err
			}

			cores = append(cores, c)
		}

		return nil, fmt.Errorf(`core "%s" should contains type property`, def)
	}

	if len(cores) == 0 {
		cores = append(cores, zapcore.NewNopCore())
	}

	return cores, err
}

// options is options getter.
func (b *Bundle) options(cfg *viper.Viper) (_ []zap.Option, err error) {
	var (
		dev     = cfg.GetBool(BundleName + ".development")
		options = make([]zap.Option, 0, 8)
	)

	if dev {
		options = append(options, zap.Development())
	}

	if cfg.GetBool(BundleName + ".caller") {
		options = append(options, zap.AddCaller())
	}

	if cfg.IsSet(BundleName + ".stacktrace") {
		var level = zap.NewAtomicLevel()
		if err = level.UnmarshalText([]byte(cfg.GetString(BundleName + ".stacktrace"))); err != nil {
			return nil, err
		}

		options = append(options, zap.AddStacktrace(level))
	}

	var fConf = make(fieldsConf, 0, 4)
	if err = cfg.UnmarshalKey(BundleName+".fields", &fConf); err != nil {
		return nil, err
	}

	if len(fConf) > 0 {
		var fields = make([]zapcore.Field, 0, len(fConf))
		for _, f := range fConf {
			fields = append(fields, zap.String(f.Key, f.Value))
		}

		options = append(options, zap.Fields(fields...))
	}

	return options, nil
}
