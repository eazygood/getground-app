package infra

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/eazygood/getground-app/internal/config"
	"github.com/go-sql-driver/mysql"
	logger "github.com/sirupsen/logrus"
	gormMySql "gorm.io/driver/mysql"
)

var (
	db *gorm.DB
)

func InitDb(cfg *config.App) *gorm.DB {
	logger.Info(cfg.Database.Name, cfg.Database.Password, cfg.Database.User, cfg.Database.Host, cfg.Database.Port)
	dbConfig := mysql.Config{
		User:                 cfg.Database.User,
		Passwd:               cfg.Database.Password,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", cfg.Database.Host, cfg.Database.Port),
		DBName:               cfg.Database.Name,
		Collation:            "utf8_general_ci",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	config := gormMySql.Config{
		DSN: dbConfig.FormatDSN(),
	}

	logger.Info(dbConfig.FormatDSN())

	initDB, err := gorm.Open(gormMySql.New(config), &gorm.Config{})

	if err != nil {
		logger.Panicf("failed to init database: %v\n", err)
	}

	db = initDB

	return db
}
