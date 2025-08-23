package repository

import (
	"app/internal/entity"
	"context"
	"database/sql"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Register(ctx context.Context, user *entity.User) (int64, error) {
	const q = `
		INSERT INTO users (email, password_hash, description, created_at, updated_at)
		VALUES ($1, $2, $3, now(), now())
		RETURNING id
	`
	var id int64
	if err := r.db.QueryRowContext(ctx, q, user.Email, user.PasswordHash, user.Description).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	const q = `
		SELECT id, email, password_hash, description, created_at, updated_at
		FROM users
		WHERE id = $1
		LIMIT 1
	`
	var u entity.User
	if err := r.db.QueryRowContext(ctx, q, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Description, &u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	const q = `
		SELECT id, email, password_hash, description, created_at, updated_at
		FROM users
		WHERE email = $1
		LIMIT 1
	`
	var u entity.User
	if err := r.db.QueryRowContext(ctx, q, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Description, &u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &u, nil
}
