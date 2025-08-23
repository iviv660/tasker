package entity

import "time"

type Task struct {
	ID          int64     `json:"id"`
	OwnerID     int64     `json:"owner_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      bool      `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
