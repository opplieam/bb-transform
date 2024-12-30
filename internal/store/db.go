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
	DBCtxTimeout         = 5 * time.Second
)

// NewDB creates a new database connection and configures it with given parameters.
func NewDB() (*sql.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DBCtxTimeout)
	defer cancel()

	// Use IPV4 for AWS lambda
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
