package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

const (
	DBDriver      string = "postgres"
	DBMaxOpenConn int    = 25
	DBMaxIdleConn int    = 25
	DBMaxIdleTime string = "15m"
	DBDSNParam    string = "BUYBETTER_DEV_SUPABASE_DSN"
)

func NewDB() (*sql.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbDSN, err := getDSN(ctx)

	if err != nil {
		return nil, err
	}
	db, err := sql.Open(DBDriver, dbDSN)
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

func getDSN(ctx context.Context) (string, error) {
	// Load the default AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", err
	}

	// Create an SSM client
	client := ssm.NewFromConfig(cfg)
	resp, err := client.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(DBDSNParam),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", err
	}
	return *resp.Parameter.Value, nil
}
