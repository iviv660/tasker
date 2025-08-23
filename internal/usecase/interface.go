package usecase

import (
	"app/internal/entity"
	"context"
)

type RepoUser interface {
	Register(ctx context.Context, user entity.User) error
	Auth(ctx context.Context, username string, password string) (entity.User, error)
	GetByID(ctx context.Context, id int64) (entity.User, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
}

type RepoTask interface {
	Create(ctx context.Context, task entity.Task) error
	Update(ctx context.Context, task entity.Task) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]entity.Task, error)
}
