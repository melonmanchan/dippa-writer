package models

import (
	_ "database/sql"
	"log"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/pkg/errors"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

// Client ...
type Client struct {
	DB *sqlx.DB
}

// PerformPendingMigrations ...
func PerformPendingMigrations(path string, connectionString string) error {
	log.Println(connectionString)
	m, err := migrate.New(path, connectionString)

	if err != nil {
		return errors.Wrap(err, "Connecting to migrations failed")
	}

	err = m.Up()

	if err != nil && err.Error() != "no change" {
		return errors.Wrap(err, "Creating migrations failed")
	}

	m.Close()

	return nil
}

// ConnectToDatabase ...
func ConnectToDatabase(connectionString string) (*Client, error) {
	log.Print("Attempting to connect to " + connectionString)

	db, err := sqlx.Connect("postgres", connectionString)

	if err != nil {
		return nil, err
	}

	return &Client{DB: db}, nil
}

// Close ...
func (c *Client) Close() {
	c.DB.Close()
}
