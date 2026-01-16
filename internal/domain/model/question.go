package model

import "time"

type Question struct {
	ID           string     `json:"id"`
	ProductID    string     `json:"product_id"`
	UserID       string     `json:"user_id"`
	UserName     string     `json:"user_name"`
	Question     string     `json:"question"`
	Answer       string     `json:"answer"`
	QuestionDate time.Time  `json:"question_date"`
	AnswerDate   *time.Time `json:"answer_date,omitempty"`
	Likes        int        `json:"likes"`
}
