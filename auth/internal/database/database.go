package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"

	"github.com/abisalde/gprc-microservice/auth/pkg/ent"
	_ "github.com/lib/pq"
)

var (
	client     *ent.Client
	clientOnce sync.Once
	initErr    error
)

type Database struct {
	Client *ent.Client
	db     *sql.DB
}

func Connect() (*Database, error) {
	var (
		db       *sql.DB
		dbClient *ent.Client
	)

	clientOnce.Do(func() {
		var err error
		db, err = initDatabase()
		if err != nil {
			initErr = fmt.Errorf("🛑 Database initialization failed: %w", err)
			return
		}

		drv := entsql.OpenDB(dialect.Postgres, db)
		dbClient = ent.NewClient(ent.Driver(drv), ent.Debug(), ent.Log(log.Print))

		if err := migrate(context.Background(), dbClient); err != nil {
			initErr = fmt.Errorf("🛠️ Database migration failed: %w", err)
			_ = dbClient.Close()
			_ = db.Close()
			return
		}
	})

	if initErr != nil {
		return nil, initErr
	}

	return &Database{
		Client: dbClient,
		db:     db,
	}, nil
}

func migrate(ctx context.Context, client *ent.Client) error {

	return client.Schema.Create(
		ctx,
		schema.WithDropIndex(true),
		schema.WithDropColumn(true),
		schema.WithForeignKeys(true),
	)
}

func formatDSN() string {
	return "postgres://microservice-user:Password123@auth_db:5432/microservice?sslmode=disable"
}

func (d *Database) Close() error {
	if d.Client == nil {
		return nil
	}

	if err := d.Client.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	client = nil
	clientOnce = sync.Once{}
	return nil
}

func initDatabase() (*sql.DB, error) {

	dsn := formatDSN()

	postgresSQL, err := sql.Open(dialect.Postgres, dsn)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to open database connection: %w", err)
	}

	postgresSQL.SetMaxIdleConns(10)
	postgresSQL.SetMaxOpenConns(100)
	postgresSQL.SetConnMaxLifetime(time.Hour)
	postgresSQL.SetConnMaxIdleTime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := postgresSQL.PingContext(ctx); err != nil {
		_ = postgresSQL.Close()
		return nil, fmt.Errorf("⚙️ Database ping failed: %w", err)
	}

	return postgresSQL, nil
}

func (d *Database) HealthCheck(ctx context.Context) error {
	if d.db == nil {
		return fmt.Errorf("🎛️ Database connection not initialized")
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
	}

	if err := d.db.PingContext(ctx); err != nil {
		return fmt.Errorf("🕹️ Database ping failed: %w", err)
	}

	_, err := d.db.ExecContext(ctx, "SELECT 1")
	if err != nil {
		return fmt.Errorf("🩸 Database schema verification failed: %w", err)
	}

	return nil
}

func WithTestTransaction(ctx context.Context, client *ent.Client, fn func(tx *ent.Tx) error) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return rerr
		}
		return err
	}

	return tx.Rollback()
}
