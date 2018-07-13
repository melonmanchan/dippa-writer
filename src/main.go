package main

import (
	"fmt"

	"github.com/melonmanchan/dippa-writer/src/config"
	"github.com/melonmanchan/dippa-writer/src/models"
)

func main() {
	config := config.ParseConfig()

	err := models.PerformPendingMigrations(config.MigrationsPath, config.DatabaseURL)
	fmt.Printf("%v", err)
}
