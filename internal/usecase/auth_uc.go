package usecase

import (
	"app/internal/entity"
	"app/internal/security"
	"context"
	"database/sql"
	"errors"
)

type UserUseCase struct {
	repo RepoUser
}

func NewUserUseCase(repo RepoUser) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (u *UserUseCase) Register(ctx context.Context, email, password, description string) (int64, error) {
	existing, err := u.repo.GetByEmail(ctx, email)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	if existing != nil {
		return 0, errors.New("аккаунт с этой почтой уже существует")
	}

	hash, err := security.HashPassword(password)
	if err != nil {
		return 0, err
	}

	id, err := u.repo.Register(ctx, &entity.User{
		Email:        email,
		PasswordHash: hash,
		Description:  description,
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (u *UserUseCase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("неверный email или пароль")
		}
		return "", err
	}
	if !security.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("неверный email или пароль")
	}
	token, err := security.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return "", err
	}
	return token, nil
}
