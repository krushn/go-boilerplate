package models

import (
	"go-boilerplate/db"
	"go-boilerplate/forms"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID          uint64 `gorm:"primaryKey"`
	Name        string `gorm:"size:50" json:"name" binding:"required"`
	Email       string `gorm:"size:100" json:"email" binding:"required"`
	Password    string `json:"-" binding:"required"`
	PhoneNumber string `gorm:"size:15" json:"phone_number"`
}

// UserModel ...
type UserModel struct{}

var authModel = new(AuthModel)

// Login ...
func (m UserModel) Login(form forms.LoginForm) (user User, token Token, err error) {

	result := db.GetDB().Model(&user).Where("email", form.Email).First(&user)

	if result.Error != nil {
		return user, token, result.Error
	}

	//Compare the password form and database if match
	bytePassword := []byte(form.Password)
	byteHashedPassword := []byte(user.Password)

	err = bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)

	if err != nil {
		return user, token, err
	}

	//Generate the JWT auth token
	tokenDetails, err := authModel.CreateToken(user.ID)
	if err != nil {
		return user, token, err
	}

	saveErr := authModel.CreateAuth(user.ID, tokenDetails)
	if saveErr == nil {
		token.AccessToken = tokenDetails.AccessToken
		token.RefreshToken = tokenDetails.RefreshToken
	} else {
		log.Panic(saveErr)
	}

	return user, token, nil
}

// Register ...
func (m UserModel) Register(form forms.RegisterForm) (user User, err error) {
	getDb := db.GetDB()

	//Check if the user exists in database
	var checkUser int64
	resultIfExists := getDb.Model(&user).Where("email", form.Email).Count(&checkUser)
	if resultIfExists.Error != nil {
		return user, errors.New("something went wrong, please try again later")
	}

	if checkUser > 0 {
		return user, errors.New("email already exists")
	}

	bytePassword := []byte(form.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		return user, errors.New("something went wrong, please try again later")
	}

	//Create the user and return back the user ID
	newUser := User{Email: form.Email, Password: string(hashedPassword), Name: form.Name}

	result := getDb.Create(&newUser) // pass pointer of data to Creat

	if result.Error != nil {
		return user, errors.New("something went wrong, please try again later")
	}

	newUser.Name = form.Name
	newUser.Email = form.Email

	return newUser, err
}

// One ...
func (m UserModel) One(userID uint64) (user User, err error) {
	result := db.GetDB().Model(&User{}).Where("ID", userID).First(&user)
	return user, result.Error
}
