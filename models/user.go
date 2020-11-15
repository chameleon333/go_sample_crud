package models

import (
	"errors"
	"log"
	"time"
	"work/crypto"
	"work/db"
)

// User モデルの宣言
type User struct {
	ID        int
	Email     string `form:"email" binding:"required" gorm:"unique;not null"`
	Password  string `form:"password" binding:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func IsUserExist(email string) (bool, User) {
	var user User
	db.DB.First(&user, "email=?", email)

	if user.ID > 0 {
		return true, user
	}

	return false, user

}

func GetUser(email, password string) (*User, error) {
	var user User
	passwordEncrypt, _ := crypto.PasswordEncrypt(password)

	// レコードが見つからない場合のエラーハンドリング
	if db.DB.Where("email = ? AND password >= ?", email, passwordEncrypt).First(&user).RecordNotFound() {
		return &user, errors.New("email \"" + email + "\" already exists")
	}

	return &user, nil
}

func (user *User) CreateUser(email string, password string) error {
	passwordEncrypt, _ := crypto.PasswordEncrypt(password)

	isExist, _ := IsUserExist(email)
	if isExist {
		return errors.New("email \"" + email + "\" already exists")
	}

	if err := db.DB.Create(&User{Email: email, Password: passwordEncrypt}).Error; err != nil {
		log.Fatal(err)
	}

	return nil
}

// func (user *User) GetLoginUser(email string, password string) *User {

// 	db.First(&user, "email =? AND password = ?", user.Email, user.Password)

// 	return &User{ID: user.ID, Email: user.Email, Password: user.Password}
// }
