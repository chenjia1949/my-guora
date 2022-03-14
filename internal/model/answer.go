package model

import (
	"log"

	"my-guora/internal/database"

	"github.com/jinzhu/gorm"
)

// Answer struct
type Answer struct {
	GORMBase
	Content          string      `json:"content" gorm:"type:varchar(4000)"`
	Type             int         `json:"type"`
	Question         Question    `json:"question" gorm:"ForeignKey:QuestionID"`
	QuestionID       int         `json:"questionID"`
	AnswerProfile    Profile     `json:"answerProfile" gorm:"ForeignKey:AnswerProfileID"`
	AnswerProfileID  int         `json:"answerProfileID"`
	Comments         []Comment   `json:"-"`
	CommentsCounts   int         `json:"commentsCounts"`
	Supporters       []Supporter `json:"-"`
	SupportersCounts int         `json:"supportersCounts"`
	Supported        bool        `json:"supported" gorm:"-"`
}

// Get func
func (a *Answer) Get() (answer Answer, err error) {

	if err = database.DB.Where(&a).Preload("Question").Preload("AnswerProfile").First(&answer).Error; err != nil {
		log.Print(err)
	}

	return
}

// Create func
func (a *Answer) Create() (ra int64, err error) {

	if err = database.DB.Create(&a).Error; err != nil {
		ra = -1
		log.Print(err)
	} else {
		ra = 1
	}

	return
}

// Update func
func (a *Answer) Update() (ra int64, err error) {

	if err = database.DB.Model(&a).Updates(a).Error; err != nil {
		ra = -1
		log.Print(err)
	} else {
		ra = 1
	}
	return
}

// Delete func
func (a *Answer) Delete() (ra int64, err error) {

	if err = database.DB.Where(&a).First(&a).Delete(&a).Error; err != nil {
		ra = -1
		log.Print(err)
	} else {
		ra = 1
	}
	return
}

// GetList func
func (a *Answer) GetList(limit int, offset int) (answers []Answer, err error) {

	if err = database.DB.Offset(offset).Limit(limit).Preload("Question").Preload("AnswerProfile").Find(&answers, a).Error; err != nil {
		log.Print(err)
	}

	return
}

// GetOrderList func
func (a *Answer) GetOrderList(limit int, offset int, order string) (answers []Answer, err error) {

	if err = database.DB.Offset(offset).Limit(limit).Preload("Question").Preload("AnswerProfile").Order(order).Find(&answers, a).Error; err != nil {
		log.Print(err)
	}

	return
}

// GetCounts func
func (a *Answer) GetCounts() (counts int, err error) {

	if err = database.DB.Model(&Answer{}).Where(&a).Count(&counts).Error; err != nil {
		log.Print(err)
	}
	return

}

// AfterCreate func
func (a *Answer) AfterCreate(tx *gorm.DB) (err error) {

	var q Question
	q.ID = a.QuestionID

	if err = tx.Model(&q).UpdateColumn("answers_counts", gorm.Expr("answers_counts + ?", 1)).Error; err != nil {
		log.Print(err)
	}

	return
}

// AfterDelete func
func (a *Answer) AfterDelete(tx *gorm.DB) (err error) {

	var q Question
	q.ID = a.QuestionID

	if err = tx.Model(&q).UpdateColumn("answers_counts", gorm.Expr("answers_counts - ?", 1)).Error; err != nil {
		log.Print(err)
	}

	return
}
