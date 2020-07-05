package postgres

import (
	"context"
	"fmt"
	pgx "github.com/jackc/pgx/v4/pgxpool"
	"regexp"
)

// pgx connection pool
var pool *pgx.Pool

// prepares the db for operation
func dbInit(host string, db string, user string, password string) error {
	dbUrl := fmt.Sprintf("postgres://%v:%v@%v:5432/%v", user, password, host, db)

	// make a new connection pool
	var err error
	pool, err = pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		return err
	}

	// execute the 'startup' stmt
	_, err = pool.Exec(context.Background(), stmtStartup)
	if err != nil {
		// check if its an actual error or just "schema already exists"
		if matched, _ := regexp.Match(".*\\(SQLSTATE 42P06\\)", []byte(err.Error())); matched {
			fmt.Println("Schema in database exists, let's assume it's correct...")
			return nil
		}
		return err
	}

	fmt.Println("Created new schema on Postgres server.")
	return nil
}
