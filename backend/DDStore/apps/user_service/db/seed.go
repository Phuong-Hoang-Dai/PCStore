package db

import (
	"context"
	"log"

	"github.com/Phuong-Hoang-Dai/DDStore/app/user_service/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedData(db *pgxpool.Pool) error {
	ctx := context.Background()

	var seeded bool
	if err := db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM users)`).Scan(&seeded); err != nil {
		return err
	}
	if seeded {
		log.Println("Data is seeded")
		return nil
	}

	user := []model.User{
		{Name: "system", Email: "system@system.com", Role: "system"},
		{Name: "admin1", Email: "admin@admin.com", Role: "admin"},
	}

	pw := []string{"systempw", "admin01"}

	for i := range user {
		user[i].SetHashPassword([]byte(pw[i]))
	}

	_, err := db.CopyFrom(ctx,
		pgx.Identifier{"users"},
		[]string{"name", "email", "password", "role"},
		pgx.CopyFromSlice(len(user), func(i int) ([]any, error) {
			return []any{user[i].Name, user[i].Email, user[i].HashedPassword, user[i].Role}, nil
		}),
	)
	return err
}
