package model

import (
	"errors"
	"log"

	"my-guora/internal/database"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User struct
type User struct {
	GORMBase
	Mail       string  `json:"mail" gorm:"type:varchar(100);unique_index"`
	Password   string  `json:"password"`
	Authorized int     `json:"authorized"`
	Type       int     `json:"type"`
	Profile    Profile `json:"profile" gorm:"ForeignKey:ProfileID"`
	ProfileID  int     `json:"profileID"`
}

// Get func
func (u *User) Get() (user User, err error) {

	if err = database.DB.Where(&u).Preload("Profile").First(&user).Error; err != nil {
		log.Print(err)
	}

	return
}

// Create func
func (u *User) Create() (ra int64, err error) {

	if err = database.DB.Create(&u).Error; err != nil {
		ra = -1
		log.Print(err)
	} else {
		ra = 1
	}

	return
}

// Update func
func (u *User) Update() (ra int64, err error) {

	if err = database.DB.Model(&u).Updates(u).Error; err != nil {
		ra = -1
		log.Print(err)
	} else {
		ra = 1
	}
	return
}

// Delete func
func (u *User) Delete() (ra int64, err error) {
	if err = database.DB.Delete(&u).Error; err != nil {
		ra = -1
		log.Print(err)
	} else {
		ra = 1
	}

	return
}

// GetList func
func (u *User) GetList(limit int, offset int) (users []User, err error) {

	if err = database.DB.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		log.Print(err)
	}

	return
}

// GetCounts func
func (u *User) GetCounts() (counts int, err error) {

	if err = database.DB.Model(&User{}).Count(&counts).Error; err != nil {
		log.Print(err)
	}

	return
}

// BeforeDelete func
func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
	if u.ID == 1 {
		err = errors.New("Can Not Remove Admin")
	}
	return
}

// AfterCreate func
func (u *User) AfterCreate(tx *gorm.DB) (err error) {

	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Print(err)
		return
	}

	if err = tx.Model(&u).UpdateColumn("password", string(bytes)).Error; err != nil {
		log.Print(err)
	}

	return
}
