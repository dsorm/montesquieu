package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/david-sorm/montesquieu/store"
	pgx "github.com/jackc/pgx/v4/pgxpool"
	"regexp"
	"time"
)

// pgx connection pool
var pool *pgx.Pool

// this is global context for this store, if this store is shutting down, the
// context is shutting down too
// it's used as a parent to all other contexts
var ctx context.Context
var ctxCancelFunc context.CancelFunc

// Init implements Store's Init function
func (p *Store) Init(f func(), cfg store.StoreConfig) error {
	p.ArticlesPerIndexPage = cfg.ArticlesPerIndexPage
	ctx, ctxCancelFunc = context.WithCancel(context.Background())

	err := dbInit(cfg.Host, cfg.Database, cfg.Username, cfg.Password, cfg.Port)

	if err != nil {
		return err
	}
	return nil
}

// Close implements Store's Close (future) function
func (p *Store) Close() {
	ctxCancelFunc()
}

// prepares the db for operation
func dbInit(host string, db string, user string, password string, port string) error {

	// make a new connection pool
	// TODO connection timeout
	var err error
	c, _ := context.WithCancel(ctx)

	// use the default port if it's undefined by the user
	if port == "" {
		port = "5432"
	}
	dsn := fmt.Sprintf("user=%v password=%v host=%v dbname=%v port=%v pool_max_conns=20",
		user, password, host, db, port)
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		panic(err)
	}

	err = errors.New("")

	// up to 30 seconds timeout for postgres
	maxCount := 30
	for count := 1; count <= maxCount; count++ {
		fmt.Println()
		connectionContext, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()
		pool, err = pgx.ConnectConfig(connectionContext, config)
		if err != nil {
			fmt.Printf("\rConnecting to postgres... (%v/%v)", count, maxCount)
			count++
		} else {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		return errors.New("Postgres connection timeout exceeded.")
	}

	// execute the 'startup' stmt
	_, err = pool.Exec(returnConnectionCtx(), stmtStartup)
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

// returns context for every connection
func returnConnectionCtx() context.Context {
	r, _ := context.WithTimeout(ctx, 5*time.Second)
	return r
}
