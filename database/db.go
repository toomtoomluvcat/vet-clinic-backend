package database

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() (err error) {
	err = godotenv.Load("../.env")
	if err!=nil{
		return fmt.Errorf("Fail to load .env file: %v",err)
	}
	Host_Name := os.Getenv("Hostname")
	Port :=os.Getenv("Port")
	User_Name :=os.Getenv("Username")
	DB_Name := os.Getenv("Database")
	Password :=os.Getenv("Password")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require search_path=public", Host_Name, User_Name, Password, DB_Name, Port)


	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err!=nil{
		return fmt.Errorf("Fail to Connect postgres: %v",err)
	}
	// err =DB.Debug().AutoMigrate(&model.Owner{}, &model.Pet{}, &model.Transaction{})
	// if err!=nil{
	// 	return fmt.Errorf("Fail to create table:%v",err)
	// }
	return nil
}

