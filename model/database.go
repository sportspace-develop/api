package model

import (
	"fmt"
	"sport-space-api/logger"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	Cfg  DSN
	DB   *gorm.DB
	log  *logger.Logger
	once sync.Once
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

	// db.Migrator()
	db.AutoMigrate(&User{})
	db.AutoMigrate(&UserSession{})
	db.AutoMigrate(&UserAuthCode{})

	db.AutoMigrate(&Team{})
	db.AutoMigrate(&TeamInvite{})

	db.AutoMigrate(&DGame{})
	db.AutoMigrate(&Player{})
	db.AutoMigrate(&Tournament{})
	db.AutoMigrate(&TournamentApplication{})
	db.AutoMigrate(&TournamentApplicationPlayer{})

	db.AutoMigrate(&Organization{})
}

func Connect() (*gorm.DB, error) {
	// var err error
	// if DB == nil {
	// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", Cfg.Username(), Cfg.Password(), Cfg.Host(), Cfg.Port(), Cfg.DBName())
	// 	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// 	if err != nil {
	// 		log.ERROR(err.Error())
	// 	}
	// }
	// return DB, err

	once.Do(func() {
		var err error
		dsn := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=disable", Cfg.Host(), Cfg.Port(), Cfg.Username(), Cfg.Password(), Cfg.DBName())
		if DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
			log.ERROR(err.Error())
		}
		// sqlDB, err := DB.DB()
		// sqlDB.SetMaxOpenConns(20)
		// sqlDB.SetMaxIdleConns(0)
		// sqlDB.SetConnMaxLifetime(time.Nanosecond)
	})
	return DB, nil
}
