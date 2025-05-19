package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/redat00/seacrate/internal/config"
)

type DatabaseEngine struct {
	c *pgx.Conn
}

func (d *DatabaseEngine) Init() error {
	// Create the folder table
	var createTableFoldersRequest = `
		CREATE TABLE IF NOT EXISTS folders (
			id SERIAL PRIMARY KEY,
			path TEXT NOT NULL UNIQUE,
			parent_folder TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
	`
	_, err := d.c.Exec(context.Background(), createTableFoldersRequest)
	if err != nil {
		return err
	}

	// Create the initial root folder
	_, err = d.c.Exec(context.Background(), "INSERT INTO folders (path) VALUES ($1) ON CONFLICT DO NOTHING", "/")
	if err != nil {
		return err
	}

	// Create the secrets table
	var createTableSecretsRequest = `
		CREATE TABLE IF NOT EXISTS secrets (
			key TEXT NOT NULL,
			value TEXT NOT NULL,
			folder TEXT REFERENCES folders(path) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			PRIMARY KEY (key, folder)
		);
	`
	_, err = d.c.Exec(context.Background(), createTableSecretsRequest)
	if err != nil {
		return err
	}

	// Create the meta table
	var createTableMetaRequest = `
		CREATE TABLE IF NOT EXISTS meta (
			key text PRIMARY KEY,
			value text NOT NULL
		);
	`
	_, err = d.c.Exec(context.Background(), createTableMetaRequest)
	if err != nil {
		return err
	}

	return nil
}

func NewDatabaseEngine(config config.DatabaseConfiguration) (DatabaseEngine, error) {
	var databaseEngine DatabaseEngine
	var err error

	connectionUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", config.Username, config.Password, config.Host, config.Port, config.Database)
	conn, err := pgx.Connect(context.Background(), connectionUrl)
	if err != nil {
		panic(err)
	}

	databaseEngine.c = conn

	err = databaseEngine.Init()
	if err != nil {
		return databaseEngine, err
	}

	return databaseEngine, err
}
