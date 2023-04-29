package models

import "auth-service/config"

func AutoMigrateModels() {
	config.DB.AutoMigrate(&Users{})
}
