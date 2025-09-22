package db

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Open(dataSource string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dataSource)
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(context.Background()); err != nil {
		return nil, err
	}
	return db, nil
}
