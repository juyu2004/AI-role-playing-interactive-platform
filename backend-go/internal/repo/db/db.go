package db

import (
	"context"
	"database/sql"
	"log"
	"net/url"

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
	log.Printf("DB connected: %s", maskDSN(dataSource))
	return db, nil
}

func maskDSN(dsn string) string {
	u, err := url.Parse(dsn)
	if err != nil {
		return "<invalid dsn>"
	}
	if u.User != nil {
		username := u.User.Username()
		u.User = url.UserPassword(username, "***")
	}
	return u.Redacted()
}
