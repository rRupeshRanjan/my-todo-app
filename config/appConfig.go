package config

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"my-todo-app/domain"
	"os"
)

var (
	Port               string
	AppLogger          *zap.Logger
	SqlDriver          string
	DataSourceName     string
	fiberLogFormat     string
	fiberLogTimeFormat string
	corsAllowOrigins   string
	corsAllowHeaders   string
	accessLogFile      *os.File
)

func init() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")

	// Look for config file in one of below paths
	viper.AddConfigPath(".")
	viper.AddConfigPath("..") // for tests only

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err == nil {
		Port = viper.GetString(domain.AppServerPort)
		SqlDriver = viper.GetString(domain.SqlDriver)
		DataSourceName = viper.GetString(domain.SqlDatabaseName)
		fiberLogFormat = viper.GetString(domain.FiberLogFormat)
		fiberLogTimeFormat = viper.GetString(domain.FiberLogTimeFormat)
		corsAllowOrigins = viper.GetString(domain.CorsAllowedOrigin)
		corsAllowHeaders = viper.GetString(domain.CorsAllowedHeaders)
		AppLogger = getLogger(domain.AppLogLocation)
		accessLogFile = getFile(domain.AppAccessLogLocation)
	} else {
		log.Panic(fmt.Sprintf("Unable to read config, program will exit now. Error: %s", err.Error()))
	}
}

func getLogger(filepath string) *zap.Logger {
	file := getFile(filepath)
	return zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(file),
			zap.InfoLevel,
		),
	)
}

func getFile(filepath string) *os.File {
	AccessLogFile, err := os.OpenFile(viper.GetString(filepath), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Panic(fmt.Sprintf("Error opening log file: %s with error: %s",
			viper.GetString(filepath), err.Error()))
	}

	return AccessLogFile
}

func GetCors() fiber.Handler {
	return cors.New(
		cors.Config{
			AllowHeaders: corsAllowHeaders,
			AllowOrigins: corsAllowOrigins,
		})
}

func GetFiberLogger() fiber.Handler {
	return logger.New(
		logger.Config{
			Format:     fiberLogFormat,
			TimeFormat: fiberLogTimeFormat,
			Output:     accessLogFile,
		})
}
