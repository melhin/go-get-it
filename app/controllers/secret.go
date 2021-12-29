package controllers

import (
	"go-get-it/app/auth"
	"net/http"

	"gorm.io/gorm"
)

func SecretView(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	userDetails, err := auth.ExtractTokenDetails(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "extraction of token failed")
		return
	}
	RespondJSON(w, http.StatusOK, userDetails)
}
