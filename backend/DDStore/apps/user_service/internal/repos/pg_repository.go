package repos

import (
	"context"

	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/internal/model"
	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/internal/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgUserRepo struct {
	DB *pgxpool.Pool
}

func NewPgUserRepo(db *pgxpool.Pool) service.UserRepos {
	return pgUserRepo{
		DB: db,
	}
}

const userColumns = `id, name, email, password, role, created_at, updated_at`

func (m pgUserRepo) CreateUser(data *model.User, ctx context.Context) (int, error) {
	err := m.DB.QueryRow(ctx,
		`INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`,
		data.Name, data.Email, data.HashedPassword, data.Role,
	).Scan(&data.Id, &data.CreatedAt, &data.UpdatedAt)

	return data.Id, err
}

// Empty fields are left unchanged.
func (m pgUserRepo) UpdateUser(data model.User, ctx context.Context) error {
	result, err := m.DB.Exec(ctx,
		`UPDATE users SET name = COALESCE(NULLIF($2, ''), name), email = COALESCE(NULLIF($3, ''), email),
		updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`,
		data.Id, data.Name, data.Email)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (m pgUserRepo) DeleteUser(id int, ctx context.Context) error {
	result, err := m.DB.Exec(ctx,
		`UPDATE users SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (m pgUserRepo) GetUserById(id int, data *model.User, ctx context.Context) error {
	return scanUser(m.DB.QueryRow(ctx,
		`SELECT `+userColumns+` FROM users WHERE id = $1 AND deleted_at IS NULL`, id), data)
}

func (m pgUserRepo) GetUserByName(name string, data *model.User, ctx context.Context) error {
	return scanUser(m.DB.QueryRow(ctx,
		`SELECT `+userColumns+` FROM users WHERE name = $1 AND deleted_at IS NULL ORDER BY id LIMIT 1`, name), data)
}

// LIMIT NULLIF($1, 0): a limit of 0 means no limit (LIMIT NULL).
func (m pgUserRepo) GetUsers(p model.Paging, data *[]model.User, ctx context.Context) error {
	rows, err := m.DB.Query(ctx,
		`SELECT `+userColumns+` FROM users WHERE deleted_at IS NULL
		ORDER BY id LIMIT NULLIF($1, 0) OFFSET $2`,
		p.Limit, p.Offset)
	if err != nil {
		return err
	}

	*data, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (u model.User, err error) {
		err = scanUser(row, &u)
		return u, err
	})
	return err
}

func scanUser(row pgx.Row, u *model.User) error {
	return row.Scan(&u.Id, &u.Name, &u.Email, &u.HashedPassword, &u.Role, &u.CreatedAt, &u.UpdatedAt)
}
