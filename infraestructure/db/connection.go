package db

import (
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	instance *gorm.DB
	once     sync.Once
)

func InitConnection() *gorm.DB {
	once.Do(func() {
		var err error
		dns := os.Getenv("DATABASE_URL")
		instance, err = gorm.Open(postgres.Open(dns), &gorm.Config{})
		if err != nil {
			log.Fatalf("ERROR INIT CONNECTION: %v", err)
		}
		dbInstance, err := instance.DB()
		if err != nil {
			log.Fatalf("ERROR GET DB INSTANCE: %v", err)
		}
		dbInstance.SetMaxIdleConns(5)            // number of ineactive connections allowed
		dbInstance.SetMaxOpenConns(10)           // number of open connections allowed
		dbInstance.SetConnMaxLifetime(time.Hour) // time to live of a connection
	})
	return instance
}

func CloseConnection() {
	dbInstance, err := instance.DB()
	if err != nil {
		log.Fatalf("ERROR CLOSE CONNECTION: %v", err)
	}
	dbInstance.Close()
}
