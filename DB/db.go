package db

import (
	"fmt"
	"os"
	models "user_admin/Models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var UserList []models.User

func Init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error in loading .env file")
	}

	DB, err = gorm.Open(postgres.Open(os.Getenv("DB")), &gorm.Config{})
	if err != nil {
		fmt.Println("database is not loaded")
		return
	}
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		fmt.Println("Error in automigrating", err)
		return
	}
}
