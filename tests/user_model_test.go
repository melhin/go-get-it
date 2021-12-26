package tests

import (
	"testing"

	"gopkg.in/go-playground/assert.v1"
	"gorm.io/gorm"

	"go-get-it/app/models"
)

func createMultipleUsers(t *testing.T, db *gorm.DB) {
	newUsers := []models.User{
		{
			ID:       1,
			Name:     "test_name1",
			Username: "test_username1",
			Password: "password1",
		},
		{
			ID:       2,
			Name:     "test_name2",
			Username: "test_username2",
			Password: "password2",
		},
	}
	for _, newUser := range newUsers {

		_, err := newUser.CreateUser(db)
		if err != nil {
			t.Errorf("Error saving users: %v\n", err)
			return
		}
	}

}

func TestSaveUser(t *testing.T) {

	db := Connect()
	tx := db.Begin()
	defer tx.Rollback()
	defer CloseDB(db)
	newUser := models.User{
		ID:       1,
		Name:     "test name",
		Username: "test username",
		Password: "password",
	}

	savedUser, err := newUser.CreateUser(tx)
	if err != nil {
		t.Errorf("Error while saving a user: %v\n", err)
		return
	}
	assert.Equal(t, newUser.ID, savedUser.ID)
	assert.Equal(t, newUser.Name, savedUser.Name)
	assert.Equal(t, newUser.Username, savedUser.Username)
}

func TestUserUtils(t *testing.T) {

	db := Connect()
	tx := db.Begin()
	defer tx.Rollback()
	defer CloseDB(db)
	createMultipleUsers(t, tx)

	user := models.User{}
	users, err := user.FindAllUsers(tx)
	if err != nil {
		t.Errorf("Error while retrieving users: %v\n", err)
		return
	}
	// 2 users
	assert.Equal(t, len(*users), 2)
	assert.Equal(t, (*users)[0].ID, uint32(1))
	assert.Equal(t, (*users)[1].ID, uint32(2))

	first_user, err := user.FindUserByID(tx, uint32(1))
	if err != nil {
		t.Errorf("Error while retrieving user with ID: %v\n", err)
		return
	}
	assert.Equal(t, first_user.Name, "test_name1")
	assert.Equal(t, models.VerifyPassword(first_user.Password, "password1"), nil)

	no_users, err := user.FindUserByID(tx, uint32(1000))
	assert.Equal(t, no_users, &models.User{})
	assert.Equal(t, err.Error(), "record not found")

}
