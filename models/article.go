package models

import "time"

// Article reprsent the article model
type Article struct {
	ID        int64     `json:"id"`
	Slug      string    `json:"slug"`
	Title     int64     `json:"title"`
	Author    Author    `json:"author"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
