package config

import (
	"gin-full-api/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open("mysql", "root:@(localhost)/gin_fullapi?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}

	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Artichel{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

	DB.Model(&models.User{}).Related(&models.Artichel{})
}
