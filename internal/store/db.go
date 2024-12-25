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

// NewDB establishes a new connection to the Postgres SQL database.
// It retrieves the Data Source Name (DSN) from AWS Systems Manager Parameter Store,
// opens a database connection using the retrieved DSN, and configures the connection pool.
// It returns a pointer to the sql.DB object representing the database connection pool
// or an error if any step fails.
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

// getDSN retrieves the Database Source Name (DSN) from AWS Systems Manager Parameter Store.
// It uses the default AWS configuration and creates an SSM client to fetch the parameter
// specified by DBDSNParam. The parameter is decrypted before being returned.
// It returns the DSN string or an error if the retrieval fails.
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
