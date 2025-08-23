package repository

import (
	"app/internal/entity"
	"context"
	"database/sql"
)

type TaskRepo struct {
	db *sql.DB
}

func NewTaskRepo(db *sql.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

// Create — создаёт задачу и возвращает заполненную сущность
func (r *TaskRepo) Create(ctx context.Context, task *entity.Task) (*entity.Task, error) {
	const query = `
		INSERT INTO tasks (owner_id, title, description, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, now(), now())
		RETURNING id, created_at, updated_at
	`

	if err := r.db.QueryRowContext(ctx, query,
		task.OwnerID,
		task.Title,
		task.Description,
		task.Status,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt); err != nil {
		return nil, err
	}
	return task, nil
}

// Update — обновляет задачу владельца и возвращает актуальную версию
func (r *TaskRepo) Update(ctx context.Context, task *entity.Task) (*entity.Task, error) {
	const query = `
		UPDATE tasks
		SET title = $1,
		    description = $2,
		    status = $3,
		    updated_at = now()
		WHERE id = $4 AND owner_id = $5
		RETURNING created_at, updated_at
	`

	if err := r.db.QueryRowContext(ctx, query,
		task.Title,
		task.Description,
		task.Status,
		task.ID,
		task.OwnerID,
	).Scan(&task.CreatedAt, &task.UpdatedAt); err != nil {
		return nil, err
	}
	return task, nil
}

// Delete — удаляет задачу только у её владельца
func (r *TaskRepo) Delete(ctx context.Context, id int64, ownerID int64) error {
	const query = `DELETE FROM tasks WHERE id = $1 AND owner_id = $2`
	res, err := r.db.ExecContext(ctx, query, id, ownerID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// GetByID — достаёт задачу по id владельца
func (r *TaskRepo) GetByID(ctx context.Context, id int64, ownerID int64) (*entity.Task, error) {
	const query = `
		SELECT id, owner_id, title, description, status, created_at, updated_at
		FROM tasks
		WHERE id = $1 AND owner_id = $2
	`
	var t entity.Task
	if err := r.db.QueryRowContext(ctx, query, id, ownerID).Scan(
		&t.ID, &t.OwnerID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &t, nil
}

// List — список задач конкретного пользователя
func (r *TaskRepo) List(ctx context.Context, ownerID int64) ([]*entity.Task, error) {
	const query = `
		SELECT id, owner_id, title, description, status, created_at, updated_at
		FROM tasks
		WHERE owner_id = $1
		ORDER BY id DESC
	`
	rows, err := r.db.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*entity.Task
	for rows.Next() {
		var t entity.Task
		if err := rows.Scan(&t.ID, &t.OwnerID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}
