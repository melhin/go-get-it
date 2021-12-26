package app

import (
	"fmt"
	"log"
	"net/http"

	"go-get-it/app/models"
	"go-get-it/config"

	"github.com/gorilla/mux"

	postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// App has router and db instances
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// Initializes the app with predefined configuration
func (a *App) Initialize(config *config.Config, logger *log.Logger) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d",
		config.DB.Host, config.DB.Username, config.DB.Password, config.DB.Name, config.DB.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Could not connect database")
	} else {
		fmt.Printf("database %s connected successfully \n", dsn)
	}

	a.DB = models.DBMigrate(db)
	a.Router = mux.NewRouter()
	a.InitializeRoutes()
}

// Run the app on it's router
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
