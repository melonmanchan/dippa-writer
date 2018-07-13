package config

import "os"

// Config ...
type Config struct {
	DatabaseURL    string
	MigrationsPath string
}

// postgres://mattes:secret@localhost:5432/database
var defaultConf = Config{
	DatabaseURL:    "postgres://mattij@localhost:5432/dippa?sslmode=disable",
	MigrationsPath: "file://migrations",
}

func ParseConfig() Config {
	cfg := defaultConf

	path := os.Getenv("MIGRATIONS_PATH")
	db_path := os.Getenv("DATABASE_URL")

	if db_path != "" {
		cfg.DatabaseURL = db_path
	}

	if path != "" {
		cfg.MigrationsPath = path
	}

	return cfg
}
