package infra

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/eazygood/getground-app/internal/config"
	"github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	gormMySql "gorm.io/driver/mysql"
)

var (
	db *gorm.DB
)

func InitDb(cfg *config.App) *gorm.DB {
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

	initDB, err := gorm.Open(gormMySql.New(config), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Panicf("failed to init database: %v\n", err)
	}

	db = initDB

	return db
}
