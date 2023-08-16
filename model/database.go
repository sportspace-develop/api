package model

import (
	"fmt"
	"sport-space-api/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Cfg DSN
	DB  *gorm.DB
	log *logger.Logger
)

type DSN interface {
	Host() string
	Port() string
	Username() string
	Password() string
	DBName() string
}

func Init(cfg DSN) {
	log = logger.New("database")
	log.INFO("init database")
	Cfg = cfg
	db, err := Connect()
	if err != nil {
		panic(err)
	}

	// defer db.Close()

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Session{})
	db.AutoMigrate(&UserAuthCode{})
}

func Connect() (*gorm.DB, error) {
	var err error
	if DB == nil {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", Cfg.Username(), Cfg.Password(), Cfg.Host(), Cfg.Port(), Cfg.DBName())
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.ERROR(err.Error())
		}
	}
	return DB, err
}
