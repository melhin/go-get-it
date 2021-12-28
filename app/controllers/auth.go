package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"go-get-it/app/auth"
	"go-get-it/app/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		RespondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		RespondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	user.Prepare()

	err = user.Validate("login")
	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	token, err := SignIn(db, user.Username, user.Password)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid username or password")
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"token": token})
}

func SignIn(db *gorm.DB, username string, password string) (string, error) {

	var err error

	user := models.User{}

	err = db.Model(models.User{}).Where("username = ?", username).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(user.ID)
}

func SignUp(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		RespondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		RespondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	user.Prepare()
	exist, err := models.UserExist(db, user.Username)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if exist {
		RespondError(w, http.StatusBadRequest, "Username exist")
		return
	}
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Crashed before save")
		return
	}
	_, err = user.CreateUser(db)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Saving user fails")
		return
	}
	RespondJSON(w, http.StatusCreated, nil)
}
