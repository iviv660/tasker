package entity

import "time"

type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"` // скрываем из JSON
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Tasks        []Task    `json:"tasks,omitempty"`
}
