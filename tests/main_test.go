package tests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"go-get-it/app/models"
	"go-get-it/config"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {

	config := config.GetConfig()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d",
		config.DB.Host, config.DB.Username, config.DB.Password, config.DB.Name, config.DB.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Could not connect database")
	}
	return db

}

func ConnectToDefaultDb(config *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s port=%d dbname=%s",
		config.DB.Host, config.DB.Username, config.DB.Password, config.DB.Port, "postgres",
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil

}

func CloseDB(db *gorm.DB) {

	sql, _ := db.DB()
	_ = sql.Close()
}

func createDB(config *config.Config) error {
	log.Printf("Creating Database %s", config.DB.Name)
	db, err := ConnectToDefaultDb(config)
	if err != nil {
		return err
	}
	defer CloseDB(db)
	log.Printf("Clearing Database %s if exist", config.DB.Name)
	err = destroyDB(config)
	if err != nil {
		return err
	}

	stmt := fmt.Sprintf("CREATE DATABASE %s;", config.DB.Name)
	if rs := db.Exec(stmt); rs.Error != nil {
		return rs.Error
	}
	return nil
}

func destroyDB(config *config.Config) error {
	log.Printf("Destroying Database %s", config.DB.Name)
	db, err := ConnectToDefaultDb(config)
	defer CloseDB(db)
	if err != nil {
		return err
	}
	stmt := fmt.Sprintf("Drop DATABASE IF EXISTS %s;", config.DB.Name)
	if rs := db.Exec(stmt); rs.Error != nil {
		return rs.Error
	}
	return nil
}

func TestMain(m *testing.M) {
	log.Print("Starting Tests")
	err := godotenv.Load(os.ExpandEnv(".test_env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	config := config.GetConfig()

	createDB(config)

	db := Connect()
	models.DBMigrate(db)
	CloseDB(db)

	code := m.Run()
	destroyDB(config)
	os.Exit(code)
	log.Print("Ending Tests")
}
