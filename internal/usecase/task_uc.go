package usecase

import (
	"app/internal/entity"
	"context"
	"database/sql"
	"golang.org/x/text/date"
)

type TaskUseCase struct {
	repo RepoTask
}

func NewTaskUseCase(repo RepoTask) *TaskUseCase {
	return &TaskUseCase{repo: repo}
}

func (t *TaskUseCase) Create(ctx context.Context, userID int64, title, description string) (int64, error) {
	task, err := t.repo.Create(
		ctx,
		&entity.Task{
			OwnerID:     userID,
			Title:       title,
			Description: description,
			Status:      true,
		})
}
