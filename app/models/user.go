package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `json:"name"`
	Username  string    `gorm:"size:255;not null;unique" json:"username"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave(db *gorm.DB) error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.ID = 0
	u.Name = html.EscapeString(strings.TrimSpace(u.Name))
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Password = html.EscapeString(strings.TrimSpace(u.Password))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "login":
		if u.Username == "" {
			return errors.New("username is required")
		}
		if u.Password == "" {
			return errors.New("password is required")
		}
		return nil

	default:
		if u.Name == "" {
			return errors.New("name is required")
		}
		if u.Username == "" {
			return errors.New("username is required")
		}
		if u.Password == "" {
			return errors.New("password is required")
		}
		return nil
	}
}

func (u *User) CreateUser(db *gorm.DB) (*User, error) {
	err := db.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	var err error
	users := []User{}
	err = db.Model(&User{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, err
}

func (u *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error) {
	err := db.Model(User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}

func UserExist(db *gorm.DB, username string) (bool, error) {
	user := &User{}
	err := db.Model(User{}).Where("username = ?", username).First(&user).Error
	if err == nil {
		return true, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else {
		return false, err
	}
}
