package dbrepo

import (
	"database/sql"

	"github.com/rasyad91/goBookings/internal/config"
	"github.com/rasyad91/goBookings/internal/repository"
)

type postgresDBrepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

type testDBrepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBrepo{
		App: a,
		DB:  conn,
	}
}

func NewTestRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBrepo{
		App: a,
	}
}
