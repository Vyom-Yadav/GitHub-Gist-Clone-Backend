package initializers

import (
	"fmt"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB(config *Config) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		config.DBHost,
		config.DBUserName,
		config.DBUserPassword,
		config.DBName,
		config.DBPort,
	)

	var gormConfig = &gorm.Config{}
	if config.AppEnv != "production" {
		newLogger := logger.New(
			log.New(os.Stdout, "\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info, // Log level
				IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
				ParameterizedQueries:      false,       // Don't include params in the SQL log
				Colorful:                  true,        // Disable color
			},
		)
		gormConfig.Logger = newLogger
	}

	DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		zap.L().Fatal("Failed to connect to the Database", zap.Error(err))
	}
	zap.L().Info("Connected Successfully to the Database")
}
