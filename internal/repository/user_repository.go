package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-base-project/internal/database"
	"go-base-project/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	FindByID(ctx context.Context, id int64) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindAll(ctx context.Context, query model.PaginationQuery) ([]*model.User, int64, error)
	Update(ctx context.Context, id int64, req model.UpdateUserRequest) (*model.User, error)
	Delete(ctx context.Context, id int64) error

}

type userRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	query := `
		INSERT INTO users (name, email, password, role, is_active, created_at, updated_at)
		VALUES (:name, :email, :password, :role, :is_active, :created_at, :updated_at)
		RETURNING id, name, email, role, is_active, created_at, updated_at
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	rows, err := r.db.NamedQueryContext(ctx, query, user)
	if err != nil {
		return nil, fmt.Errorf("gagal create user: %w", err)
	}
	defer rows.Close()

	var created model.User
	if rows.Next() {
		if err := rows.StructScan(&created); err != nil {
			return nil, fmt.Errorf("gagal scan hasil create: %w", err)
		}
	}

	return &created, nil
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	query := `
		SELECT id, name, email, password, role, is_active, created_at, updated_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, &user, query, id)
	if err == sql.ErrNoRows {
		return nil, nil 
	}
	if err != nil {
		return nil, fmt.Errorf("gagal find user by id: %w", err)
	}

	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	query := `
		SELECT id, name, email, password, role, is_active, created_at, updated_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, &user, query, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("gagal find user by email: %w", err)
	}

	return &user, nil
}

func (r *userRepository) FindAll(ctx context.Context, query model.PaginationQuery) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	countQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`
	if query.Search != "" {
		countQuery += ` AND (name ILIKE '%' || $1 || '%' OR email ILIKE '%' || $1 || '%')`
		err := r.db.GetContext(ctx, &total, countQuery, query.Search)
		if err != nil {
			return nil, 0, fmt.Errorf("gagal count users: %w", err)
		}
	} else {
		err := r.db.GetContext(ctx, &total, countQuery)
		if err != nil {
			return nil, 0, fmt.Errorf("gagal count users: %w", err)
		}
	}

	dataQuery := `
		SELECT id, name, email, role, is_active, created_at, updated_at
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	err := r.db.SelectContext(ctx, &users, dataQuery, query.Limit, query.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("gagal get all users: %w", err)
	}

	return users, total, nil
}

func (r *userRepository) Update(ctx context.Context, id int64, req model.UpdateUserRequest) (*model.User, error) {
	query := `
		UPDATE users
		SET name = COALESCE(NULLIF($1, ''), name),
		    email = COALESCE(NULLIF($2, ''), email),
		    updated_at = $3
		WHERE id = $4 AND deleted_at IS NULL
		RETURNING id, name, email, role, is_active, created_at, updated_at
	`

	var user model.User
	err := r.db.QueryRowxContext(ctx, query, req.Name, req.Email, time.Now(), id).StructScan(&user)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("gagal update user: %w", err)
	}

	return &user, nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE users SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("gagal delete user: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user tidak ditemukan")
	}

	return nil
}
