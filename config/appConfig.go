package config

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"my-todo-app/domain"
	"os"
)

var (
	Port           string
	AppLogger      *zap.Logger
	SqlDriver      string
	DatabaseName   string
	FiberLogFormat string
	LogFile        *os.File
)

func init() {
	viper.SetConfigFile("C:\\Users\\Rupesh\\Desktop\\workspace\\golang\\my-todo-app\\config.yml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err == nil {
		Port = domain.PortSemicolon + viper.GetString(domain.AppServerPort)
		SqlDriver = viper.GetString(domain.SqlDriver)
		DatabaseName = viper.GetString(domain.SqlDatabaseName)

		FiberLogFormat = viper.GetString(domain.FiberLogFormat)
		LogFile, err = os.OpenFile(viper.GetString(domain.AppLogLocation), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

		if err != nil {
			log.Panic(fmt.Sprintf("Error opening log file: %s with error: %s",
				viper.GetString(domain.AppLogLocation), err.Error()))
		}

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(LogFile),
			zap.InfoLevel)
		AppLogger = zap.New(core)
	} else {
		log.Panic(fmt.Sprintf("Unable to read config, program will exit now. Error: %s", err.Error()))
	}
}
