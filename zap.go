package zap

import (
	"fmt"
	"os"

	"github.com/gozix/viper"
	"github.com/sarulabs/di"
	"github.com/snovichkov/zap-gelf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	// Bundle implements the glue.Bundle interface.
	Bundle struct{}

	// Logger is type alias of zap.Logger
	Logger = zap.Logger

	// loggerConf is logger configuration struct.
	loggerConf struct {
		Cores []struct {
			Addr     string
			Host     string
			Level    string
			Encoding string
		}
		Caller bool
		Fields []struct {
			Key   string
			Value string
		}
		Stacktrace  string
		Development bool
	}

	// nopSyncer is stdout wrapper.
	nopSyncer struct {
		zapcore.WriteSyncer
	}
)

// BundleName is default definition name.
const BundleName = "zap"

// NewBundle create bundle instance.
func NewBundle() *Bundle {
	return new(Bundle)
}

// Key implements the glue.Bundle interface.
func (b *Bundle) Name() string {
	return BundleName
}

// Build implements the glue.Bundle interface.
func (b *Bundle) Build(builder *di.Builder) error {
	return builder.Add(di.Def{
		Name: BundleName,
		Build: func(ctn di.Container) (_ interface{}, err error) {
			var cfg *viper.Viper
			if err = ctn.Fill(viper.BundleName, &cfg); err != nil {
				return nil, err
			}

			var conf loggerConf
			if err = cfg.UnmarshalKey(BundleName, &conf); err != nil {
				return nil, err
			}

			var cores = make([]zapcore.Core, 0, 2)
			for _, logger := range conf.Cores {
				var core zapcore.Core
				switch logger.Encoding {
				case "console", "json":
					var eConf zapcore.EncoderConfig
					if conf.Development {
						eConf = zap.NewDevelopmentEncoderConfig()
					} else {
						eConf = zap.NewProductionEncoderConfig()
					}

					var level zap.AtomicLevel
					if len(logger.Level) > 0 {
						if err = level.UnmarshalText([]byte(logger.Level)); err != nil {
							return nil, err
						}
					}

					var enc = zapcore.NewConsoleEncoder(eConf)
					if logger.Encoding == "json" {
						enc = zapcore.NewJSONEncoder(eConf)
					}

					core = zapcore.NewCore(enc, &nopSyncer{os.Stdout}, level)
				case "gelf":
					var options = make([]gelf.Option, 0, 3)
					if len(logger.Addr) > 0 {
						options = append(options, gelf.Addr(logger.Addr))
					}

					if len(logger.Level) > 0 {
						options = append(options, gelf.LevelString(logger.Level))
					}

					if len(logger.Host) == 0 {
						if logger.Host, err = os.Hostname(); err != nil {
							return nil, err
						}
					}

					options = append(
						options,
						gelf.Host(logger.Host),
					)

					if core, err = gelf.NewCore(options...); err != nil {
						return nil, err
					}
				default:
					return nil, fmt.Errorf("unknown encoding %s", logger.Encoding)
				}

				cores = append(cores, core)
			}

			if len(cores) == 0 {
				cores = append(cores, zapcore.NewNopCore())
			}

			var options = make([]zap.Option, 0, 8)
			if conf.Caller {
				options = append(options, zap.AddCaller())
			}

			if conf.Development {
				options = append(options, zap.Development())
			}

			var level zap.AtomicLevel
			if len(conf.Stacktrace) > 0 {
				if err = level.UnmarshalText([]byte(conf.Stacktrace)); err != nil {
					return nil, err
				}

				options = append(options, zap.AddStacktrace(level))
			}

			if len(conf.Fields) > 0 {
				var fields = make([]zapcore.Field, 0, len(conf.Fields))
				for _, field := range conf.Fields {
					fields = append(fields, zap.String(field.Key, field.Value))
				}

				options = append(options, zap.Fields(fields...))
			}

			return zap.New(zapcore.NewTee(cores...), options...), nil
		},
		Close: func(obj interface{}) (err error) {
			return obj.(*zap.Logger).Sync()
		},
	})
}

// DependsOn implements the glue.DependsOn interface.
func (b *Bundle) DependsOn() []string {
	return []string{viper.BundleName}
}

// Sync is override original close.
func (nopSyncer) Sync() error {
	return nil
}
