package gorm

import (
	"database/sql"
	"mjo/config"
	"mjo/util/logger"

	_ "github.com/go-sql-driver/mysql"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

const LOG_IDENTIFIER = "UTIL_GORM"

var log = logger.SetIdentifierField(LOG_IDENTIFIER)

func InitGorm(dbCon *sql.DB) (*gorm.DB, error) {
	gormMysqlInstance := gormMysql.New(gormMysql.Config{Conn: dbCon})
	var gormCon *gorm.DB
	var err error
	if config.GetConfig().ProdMode {
		gormCon, err = gorm.Open(gormMysqlInstance, &gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Silent),
		})
	} else {
		gormCon, err = gorm.Open(gormMysqlInstance, &gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Info),
		})
	}
	if err != nil {
		log.Error(err)
		panic(err)
	}
	return gormCon, err
}