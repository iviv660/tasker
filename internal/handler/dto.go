package handler

import "app/internal/entity"

// RegisterRequest ...
type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Description string `json:"description"`
}

// LoginRequest ...
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateTaskRequest ...
type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// UpdateTaskRequest ...
type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
}

type TasksResponse struct {
	Tasks []*entity.Task `json:"tasks"`
}
