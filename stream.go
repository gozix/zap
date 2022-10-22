// Copyright 2018 Sergey Novichkov. All rights reserved.
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package zap

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	// CoreFactory is zap core factory.
	CoreFactory interface {
		Name() string
		New(cfg *viper.Viper, path string) (zapcore.Core, error)
	}

	coreFactory struct{}
)

// coreFactory implements CoreFactory.
var _ CoreFactory = (*coreFactory)(nil)

func (c *coreFactory) Name() string {
	return "stream"
}

func (c *coreFactory) New(cfg *viper.Viper, path string) (_ zapcore.Core, err error) {
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
}
