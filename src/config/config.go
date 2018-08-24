package config

import "os"

// Config ...
type Config struct {
	DatabaseURL    string
	MigrationsPath string
	RabbitMQConn   string
}

// postgres://mattes:secret@localhost:5432/database
var defaultConf = Config{
	DatabaseURL:    "postgres://mattij@localhost:5432/dippa?sslmode=disable",
	MigrationsPath: "file://migrations",
	RabbitMQConn:   "amqp://guest:guest@localhost:5672/",
}

func ParseConfig() Config {
	cfg := defaultConf

	path := os.Getenv("MIGRATIONS_PATH")
	db_path := os.Getenv("DATABASE_URL")
	rabbitmqConn := os.Getenv("RABBITMQ_ADDRESS")

	if db_path != "" {
		cfg.DatabaseURL = db_path
	}

	if path != "" {
		cfg.MigrationsPath = path
	}

	if rabbitmqConn != "" {
		cfg.RabbitMQConn = rabbitmqConn
	}

	return cfg
}
