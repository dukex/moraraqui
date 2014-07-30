package models

import (
	"os"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

var DB gorm.DB

func init() {
	databaseUrl, _ := pq.ParseURL(os.Getenv("DATABASE_URL"))
	DB, _ = gorm.Open("postgres", databaseUrl)
	DB.LogMode(false)

	DB.AutoMigrate(Property{})
	DB.AutoMigrate(RealState{})

	DB.Model(Property{}).AddIndex("idx_property_url", "url")
	DB.Model(Property{}).AddUniqueIndex("idxu_property_url", "url")
}
