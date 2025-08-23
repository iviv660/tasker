package usecase

import (
	"app/internal/entity"
	"context"
)

type RepoUser interface {
	Register(ctx context.Context, user *entity.User) (int64, error)
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}

type RepoTask interface {
	Create(ctx context.Context, task *entity.Task) (*entity.Task, error)
	Update(ctx context.Context, task *entity.Task) (*entity.Task, error)
	Delete(ctx context.Context, id int64, ownerID int64) error
	GetByID(ctx context.Context, id int64, ownerID int64) (*entity.Task, error)
	List(ctx context.Context, ownerID int64) ([]*entity.Task, error)
}
