package main

import (
	"fmt"
	"log"

	"github.com/melonmanchan/dippa-writer/src/config"
	"github.com/melonmanchan/dippa-writer/src/models"
)

func main() {
	config := config.ParseConfig()

	err := models.PerformPendingMigrations(config.MigrationsPath, config.DatabaseURL)
	fmt.Printf("%v", err)

	c, err := models.ConnectToDatabase(config.DatabaseURL)

	if err != nil {
		log.Fatal(err)
	}

	user := models.WatsonResult{
		Contents: "hello world",
		UserID:   1,
		RoomID:   1,
		Keywords: []models.Keyword{
			models.Keyword{
				Contents:  "hello world",
				Sentiment: 1,
				Relevance: 1,
				Sadness:   1,
				Joy:       1,
				Fear:      1,
				Disgust:   1,
				Anger:     1,
			},
		},
	}

	err = c.CreateWatsonResult(&user)

	log.Println(err)
}
