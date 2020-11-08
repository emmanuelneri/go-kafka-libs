package logs

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

const (
	Lib         = "Confluent"
	ProjectType = "consumer"
)

func init() {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.LevelKey = "log.level"
	loggerConfig.EncoderConfig.MessageKey = "message"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	loggerProduction, err := loggerConfig.Build()
	if err != err {
		panic(errors.Wrap(err, "fail to start Logger"))
	}

	Logger = loggerProduction
}
