package main

import (
	"mjo/config"
	"mjo/repository/mysql"
	"mjo/repository/user"
	"mjo/repository/merchant"
	"mjo/repository/outlet"
	"mjo/repository/transaction"
	"mjo/util/gorm"
	"mjo/util/logger"
	"time"

	_ "github.com/go-sql-driver/mysql"
	logRotator "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

const LOG_IDENTIFIER = "APP_CLI_MIGRATE"

var log = logger.SetIdentifierField(LOG_IDENTIFIER)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "timestamp",
			logrus.FieldKeyMsg:  "message",
		},
	})
	logrus.SetLevel(logrus.InfoLevel)

	logFile, err := logRotator.New(
		"resource/log/log.%Y%m%d%H%M",
		logRotator.WithLinkName("resource/log/log"),
		logRotator.WithMaxAge(time.Duration(86400)*time.Second),
		logRotator.WithRotationTime(time.Duration(86400)*time.Second),
	)
	if err == nil {
		logrus.SetOutput(logFile)
	} else {
		log.Error(err)
		panic(err)
	}
	defer logFile.Close()
	dbCon, err := mysql.Connect(config.GetConfig().MysqlDb)
	if err != nil {
		log.Error(err)
		panic(err)
	}

	gormDb, err := gorm.InitGorm(dbCon)

	if err != nil {
		log.Error(err)
		panic(err)
	}

	user.Migrate(*gormDb)
	merchant.Migrate(*gormDb)
	outlet.Migrate(*gormDb)
	transaction.Migrate(*gormDb)

	dbCon.Close()
}