package zenkit

import (
	"context"

	stackdriver "github.com/TV4/logrus-stackdriver-formatter"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	LogTenantField  = "zing.tnt"
	LogUserField    = "zing.usr"
	LogServiceField = "zing.svc"
)

func ContextLogger(ctx context.Context) *logrus.Entry {
	return ctxlogrus.Extract(ctx)
}

func Logger(name string) *logrus.Entry {
	InitConfig(name)	// necessary for proper logger configuration
	log := logrus.New()
	if viper.GetBool(LogStackdriverConfig) {
		log.Formatter = stackdriver.NewFormatter(
			stackdriver.WithService(name),
			// TODO: stackdriver.WithVersion
		)
	}
	cfgLevel := viper.GetString(LogLevelConfig)
	level, err := logrus.ParseLevel(cfgLevel)
	if err != nil {
		log.WithFields(logrus.Fields{
			"level":         cfgLevel,
			LogServiceField: name,
		}).Error("unable to parse log level config; defaulted to INFO")
		level = logrus.InfoLevel
	}
	log.SetLevel(level)
	return log.WithFields(logrus.Fields{
		LogServiceField: name,
	})
}
