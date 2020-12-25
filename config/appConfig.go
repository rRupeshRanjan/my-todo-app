package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"my-todo-app/domain"
	"os"
)

var (
	Port string
	AppLogger *zap.Logger
	SqlDriver string
	DatabaseName string
	FiberLogFormat string
	LogFile *os.File
)

func init()  {
	viper.SetConfigFile("config.yml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err == nil {
		Port = domain.PortSemicolon + viper.GetString(domain.AppServerPort)
		SqlDriver = viper.GetString(domain.SqlDriver)
		DatabaseName = viper.GetString(domain.SqlDatabaseName)

		FiberLogFormat = viper.GetString(FiberLogFormat)
		LogFile, _ := os.OpenFile(viper.GetString(domain.AppLogLocation), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(LogFile),
			zap.InfoLevel)
		AppLogger = zap.New(core)
	}
}