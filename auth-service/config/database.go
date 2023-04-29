package config

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func ConnectToPostgres() error {
	var counts int64
	dsn := os.Getenv("DSN")
	for {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("Connected to Postgres!")
			DB = db
			return nil
		}

		log.Println("Postgres not yet ready ...")
		counts++

		if counts > 10 {
			return err
		}

		log.Println("Backing off for two seconds....")
		time.Sleep(2 * time.Second)
		continue
	}
}
