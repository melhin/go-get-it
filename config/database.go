package config

import (
	"go-get-it/seeds"
	"os"
	"strconv"
)

type Config struct {
	DB   *DBConfig
	Seed *seeds.SeedOption
}

type DBConfig struct {
	Dialect  string
	Host     string
	Port     int
	Username string
	Password string
	Name     string
	Charset  string
}

func getenvBool(key string) (bool, error) {
	s := os.Getenv(key)
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}
	return v, nil
}

func GetConfig() *Config {
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	cleardb, _ := getenvBool("CLEAR_DB")
	populatedb, _ := getenvBool("POPULATE_DB")

	return &Config{
		DB: &DBConfig{
			Dialect:  os.Getenv("DB_DIALICT"),
			Host:     os.Getenv("DB_HOST"),
			Port:     port,
			Username: os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			Charset:  os.Getenv("DB_CHARSET"),
		},
		Seed: &seeds.SeedOption{
			ClearDB:    cleardb,
			PopulateDb: populatedb,
		},
	}
}
