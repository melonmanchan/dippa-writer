package config

import "os"

// Config ...
type Config struct {
	DatabaseURL    string
	MigrationsPath string
}

var defaultConf = Config{
	DatabaseURL:    "postres://localhost:5432",
	MigrationsPath: "./migrations",
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
