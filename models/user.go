package models

import "github.com/jinzhu/gorm"

// User モデルの宣言
type User struct {
	gorm.Model
	Email    string `form:"email" binding:"required" gorm:"unique;not null"`
	Password string `form:"password" binding:"required"`
}

// func GetLoginUser(email string, password string) *User {
// 	db := gormConnect()
// 	defer db.Close()

// 	var user User
// 	db.First(&user, "email =? AND password = ?", email, password)

// 	return &User{ID: user.ID, Email: user.Email, Password: user.Password}
// }
