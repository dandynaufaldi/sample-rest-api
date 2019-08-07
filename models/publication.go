package models

import "time"

// Publication represent publication model
type Publication struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Writers   []Author  `json:"writers"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
