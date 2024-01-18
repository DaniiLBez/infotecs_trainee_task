package pgdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"infotecs_trainee_task/internal/entity"
	"infotecs_trainee_task/internal/repo"
	"infotecs_trainee_task/pkg/postgres"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) CreateUser(ctx context.Context, user entity.User) (uuid.UUID, error) {

	sql, args, _ := r.Builder.Insert("users").
		Columns("username", "password").
		Values(user.Username, user.Password).
		Suffix("RETURNING uuid").
		ToSql()

	var userUUID uuid.UUID
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&userUUID)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return uuid.Nil, repo.ErrAlreadyExist
			}
		}
		return uuid.Nil, fmt.Errorf("UserRepo.CreateUser - r.Pool.QueryRow: %v", err)
	}

	return userUUID, nil
}

// GetUserById To delete may be unused
func (r *UserRepo) GetUserById(ctx context.Context, uuid uuid.UUID) (entity.User, error) {
	sql, args, _ := r.Builder.
		Select("uuid", "username", "password", "created_at").
		From("user").
		Where("uuid = ?", uuid).
		ToSql()

	var user entity.User
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&user)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, repo.ErrNotFound
		}
		return entity.User{}, fmt.Errorf("UserRepo.GetUserById - r.Pool.QueryRow: %v", err)
	}

	return user, nil
}

// GetUserByUsername To delete may be unused
func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	sql, args, _ := r.Builder.
		Select("uuid", "username", "password", "created_at").
		From("user").
		Where("username = ?", username).
		ToSql()

	var user entity.User
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&user)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, repo.ErrNotFound
		}
		return entity.User{}, fmt.Errorf("UserRepo.GetUserByUsername - r.Pool.QueryRow: %v", err)
	}

	return user, nil
}

func (r *UserRepo) GetUserByUsernameAndPassword(ctx context.Context, username, password string) (entity.User, error) {
	sql, args, _ := r.Builder.
		Select("uuid", "username", "password", "created_at").
		From("user").
		Where(
			squirrel.And{
				squirrel.Eq{"username": username},
				squirrel.Eq{"password": password},
			}).
		ToSql()

	var user entity.User
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&user)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, repo.ErrNotFound
		}
		return entity.User{}, fmt.Errorf("UserRepo.GetUserByUsername - r.Pool.QueryRow: %v", err)
	}

	return user, nil
}
