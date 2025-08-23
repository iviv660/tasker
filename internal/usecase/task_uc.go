package usecase

import (
	"app/internal/entity"
	"context"
)

type TaskUseCase struct {
	repo RepoTask
}

func NewTaskUseCase(repo RepoTask) *TaskUseCase {
	return &TaskUseCase{repo: repo}
}

func (t *TaskUseCase) CreateTask(ctx context.Context, userID int64, title, description string) (*entity.Task, error) {
	task, err := t.repo.Create(
		ctx,
		&entity.Task{
			OwnerID:     userID,
			Title:       title,
			Description: description,
			Status:      false,
		})
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (t *TaskUseCase) UpdateTask(ctx context.Context, task *entity.Task) (*entity.Task, error) {
	return t.repo.Update(ctx, task)
}

func (t *TaskUseCase) DeleteTask(ctx context.Context, taskID, ownerID int64) error {
	return t.repo.Delete(ctx, taskID, ownerID)
}

func (t *TaskUseCase) GetTaskByID(ctx context.Context, taskID, ownerID int64) (*entity.Task, error) {
	return t.repo.GetByID(ctx, taskID, ownerID)
}

func (t *TaskUseCase) ListTasks(ctx context.Context, ownerID int64) ([]*entity.Task, error) {
	return t.repo.List(ctx, ownerID)
}
