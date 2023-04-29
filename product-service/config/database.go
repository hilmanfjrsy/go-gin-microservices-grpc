package config

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func ConnectToMySQL() error {
	var counts int64
	dsn := os.Getenv("DSN")
	for {
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("Connected to MySQL!")
			DB = db
			return nil
		}

		log.Println("MySQL not yet ready ...")
		counts++

		if counts > 20 {
			return err
		}

		log.Println("Backing off for two seconds....")
		time.Sleep(2 * time.Second)
		continue
	}
}
