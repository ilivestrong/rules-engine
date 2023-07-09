package helpers

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type (
	RulesEngineRepo interface {
		AddApprovedPhone(ctx context.Context, phone string) error
	}

	rulesEngineRepo struct {
		Conn *pgx.Conn
	}

	Config struct {
		User   string
		Pass   string
		Host   string
		DBName string
	}
)

func (rer *rulesEngineRepo) AddApprovedPhone(ctx context.Context, phone string) error {
	_, err := rer.Conn.Exec(ctx, fmt.Sprintf("INSERT INTO approved_phones VALUES('%s')", phone))
	return err
}

func NewRulesEngineRepo(ctx context.Context, config *Config) (*rulesEngineRepo, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", config.User, config.Pass, config.Host, config.DBName)
	fmt.Println(dbURL)
	c, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialise database connection: %v", err)
	}

	return &rulesEngineRepo{
		Conn: c,
	}, nil
}
