package controllers

import (
	"net/http"

	"gorm.io/gorm"
)

func Home(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, http.StatusOK, map[string]string{"message": "Welcome To This Awesome API."})
}
