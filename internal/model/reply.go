package model

import (
	"log"

	"my-guora/internal/database"

	"github.com/jinzhu/gorm"
)

// Reply struct
type Reply struct {
	GORMBase
	Content            string  `json:"content"`
	Type               int     `json:"type"`
	Comment            Comment `json:"-" gorm:"ForeignKey:CommentID"`
	CommentID          int     `json:"commentID"`
	ReplyFromProfile   Profile `json:"replyFromProfile" gorm:"ForeignKey:ReplyFromProfileID"`
	ReplyFromProfileID int     `json:"replyFromProfileID"`
	ReplyToProfile     Profile `json:"replyToProfile" gorm:"ForeignKey:ReplyToProfileID"`
	ReplyToProfileID   int     `json:"replyToProfileID"`
}

// Get func
func (r *Reply) Get() (reply Reply, err error) {

	if err = database.DB.Where(&r).Preload("ReplyFromProfile").Preload("ReplyToProfile").First(&reply).Error; err != nil {
		log.Print(err)
	}

	return
}

// Create func
func (r *Reply) Create() (ra int64, err error) {

	if err = database.DB.Create(&r).Error; err != nil {
		ra = -1
		log.Print(err)
	} else {
		ra = 1
	}

	return
}

// Update func
func (r *Reply) Update() (ra int64, err error) {

	if err = database.DB.Model(&r).Updates(r).Error; err != nil {
		ra = -1
		log.Print(err)
	} else {
		ra = 1
	}
	return
}

// Delete func
func (r *Reply) Delete() (ra int64, err error) {

	if err = database.DB.Where(&r).First(&r).Delete(&r).Error; err != nil {
		ra = -1
		log.Print(err)
	} else {
		ra = 1
	}
	return
}

// GetList func
func (r *Reply) GetList(limit int, offset int) (replies []Reply, err error) {

	if err = database.DB.Offset(offset).Limit(limit).Preload("ReplyFromProfile").Preload("ReplyToProfile").Find(&replies, r).Error; err != nil {
		log.Print(err)
	}

	return
}

// GetCounts func
func (r *Reply) GetCounts() (counts int, err error) {

	if err = database.DB.Model(&Reply{}).Where(&r).Count(&counts).Error; err != nil {
		log.Print(err)
	}

	return
}

// AfterCreate func
func (r *Reply) AfterCreate(tx *gorm.DB) (err error) {

	var co Comment
	co.ID = r.CommentID

	if err = tx.Model(&co).UpdateColumn("replies_counts", gorm.Expr("replies_counts + ?", 1)).Error; err != nil {
		log.Print(err)
	}

	return
}

// AfterDelete func
func (r *Reply) AfterDelete(tx *gorm.DB) (err error) {

	var co Comment
	co.ID = r.CommentID

	if err = tx.Model(&co).UpdateColumn("replies_counts", gorm.Expr("replies_counts - ?", 1)).Error; err != nil {
		log.Print(err)
	}

	return
}
