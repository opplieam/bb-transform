package store

import (
	"context"
	"database/sql"
	"os"
	"time"
)

const (
	DBDriver      string = "postgres"
	DBMaxOpenConn int    = 25
	DBMaxIdleConn int    = 25
	DBMaxIdleTime string = "15m"
	DBDSN         string = "BUYBETTER_DEV_SUPABASE_DSN"
)

func NewDB() (*sql.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sql.Open(DBDriver, os.Getenv(DBDSN))
	if err != nil {
		return nil, err
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(DBMaxOpenConn)
	db.SetMaxIdleConns(DBMaxIdleConn)
	duration, err := time.ParseDuration(DBMaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)
	return db, nil
}
