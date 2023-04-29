package models

import "product-service/config"

func AutoMigrateModels() {
	config.DB.AutoMigrate(&Products{})
}
