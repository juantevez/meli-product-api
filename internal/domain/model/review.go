package model

import "time"

type Review struct {
	ID           string    `json:"id"`
	ProductID    string    `json:"product_id"`
	UserID       string    `json:"user_id"`
	UserName     string    `json:"user_name"`
	Rating       int       `json:"rating"`
	Title        string    `json:"title"`
	Comment      string    `json:"comment"`
	CreatedAt    time.Time `json:"created_at"`
	HelpfulCount int       `json:"helpful_count"`
}
