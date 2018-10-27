package main

import (
	"log"
	"sync"

	"github.com/golang/protobuf/proto"
	types "github.com/melonmanchan/dippa-proto/build/go"
	"github.com/melonmanchan/dippa-writer/src/config"
	"github.com/melonmanchan/dippa-writer/src/models"
	"github.com/melonmanchan/dippa-writer/src/rabbit"
	"github.com/streadway/amqp"
)

func consumeImagesData(db *models.Client, chann <-chan amqp.Delivery, wg *sync.WaitGroup) {
	for d := range chann {

		googleFacialRecognitionRes := &types.GoogleFacialRecognition{}

		if err := proto.Unmarshal(d.Body, googleFacialRecognitionRes); err != nil {
			log.Print("Failed to parse input for google: ", err)
			break
		}

		if googleFacialRecognitionRes.Emotion == nil {
			log.Print("Emotion is nil, not proceeding")
			break
		}

		imageRes, user, room := models.GoogleProtoToGoStructs(*googleFacialRecognitionRes)

		userId, err := db.CreateUserByName(&user)

		if err != nil {
			log.Println(err)
			break
		}

		roomId, err := db.CreateRoomByName(&room)

		if err != nil {
			log.Println(err)
			break
		}

		imageRes.UserID = userId
		imageRes.RoomID = roomId

		err = db.CreateGoogleResults(&imageRes)

		if err != nil {
			log.Println(err)
		}
	}
	wg.Done()
}

func consumeTextData(db *models.Client, chann <-chan amqp.Delivery, wg *sync.WaitGroup) {
	for d := range chann {

		watsonTextRes := &types.WatsonNLP{}

		if err := proto.Unmarshal(d.Body, watsonTextRes); err != nil {
			log.Print("Failed to parse input for watson: ", err)
			break
		}

		log.Printf("%v \n", watsonTextRes)

		textRes, user, room := models.WatsonProtoToGoStructs(*watsonTextRes)

		userId, err := db.CreateUserByName(&user)

		if err != nil {
			log.Println(err)
			break
		}

		roomId, err := db.CreateRoomByName(&room)

		if err != nil {
			log.Println(err)
			break
		}

		textRes.UserID = userId
		textRes.RoomID = roomId

		err = db.CreateWatsonResult(&textRes)

		if err != nil {
			log.Println(err)
		}
	}
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	config := config.ParseConfig()

	err := models.PerformPendingMigrations(config.MigrationsPath, config.DatabaseURL)

	c, err := models.ConnectToDatabase(config.DatabaseURL)

	if err != nil {
		log.Fatalf("%s", err)
	}

	chann, err := rabbit.TryConnectToQueue(config.RabbitMQConn)

	if err != nil {
		log.Fatalf("%s", err)
	}

	defer chann.Close()

	images, err := rabbit.GetChannelToConsume(chann, "google_results")

	if err != nil {
		log.Fatalf("%s", err)
	}

	text, err := rabbit.GetChannelToConsume(chann, "ibm_text")

	if err != nil {
		log.Fatalf("%s", err)
	}

	log.Println("startup succesful")
	log.Println("has the cache updated")

	go consumeImagesData(c, images, &wg)
	go consumeTextData(c, text, &wg)

	wg.Wait()

	log.Println("All goroutines exitted")
}
