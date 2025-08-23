package usecase

import (
	"app/internal/entity"
	"context"
	"errors"
)

type UserUseCase struct {
	repo       RepoUser
	secretCode string
}

func NewUserUseCase(repo RepoUser, secretCode string) *UserUseCase {
	return &UserUseCase{repo: repo, secretCode: secretCode}
}

func (u *UserUseCase) Register(ctx context.Context, email, password string) error {
	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user != nil {
		return errors.New("аккаунт с этой почтой уже существует")
	}

}
