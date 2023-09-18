package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type Sender struct {
	Name string `json:"name"`
	Id string `json:"id"`
	Email string `json:"email"`
}

type Receiver struct {
	Email string `json:"email"`
	URL string	`json:"url"`
}

type WebhookReceiver struct {
	URL string `json:"url"`
	// TODO: good place for any security headers
}

type EmailReceiver struct {
	Email string `json:"email"`
}

type EmailPayload struct {
	Body string `json:"body"`
	Subject string `json:"subject"`
	From string `json:"from"`
}


type WebhookPayload struct {
	Type string	`json:"type"`
	Resource string `json:"resource"`
	Object interface{} 	`json:"object"`
}

type WebhookEvent struct {
	Type string	`json:"type"`
	Sender Sender	`json:"sender"`
	Receiver  WebhookReceiver	`json:"receiver"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	Payload WebhookPayload `json:"payload"`
}

type EmailEvent struct {
	Sender Sender `json:"sender"`
	Receiver EmailReceiver	`json:"receiver"`
	Payload EmailPayload	`json:"payload"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}



func main(){

	APP_ENV := os.Getenv("APP_ENV"); if APP_ENV == "DEVELOPMENT" || len(APP_ENV) == 0 {
		if err := godotenv.Load(); err != nil {
			panic("UNABLE TO LOAD ENV")
		}
	}

	MONGO_CONNECTION_STRING := os.Getenv("MONGO_CONNECTION_STRING"); if len(MONGO_CONNECTION_STRING) == 0 {
		panic("MONGO_CONNECTION_STRING NOT SET!!!")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().ApplyURI(MONGO_CONNECTION_STRING).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts); if err != nil {
		log.Panic("Unable to initialize mongo client")
	}

	log.Println("Mongo Initialized successfully")

	emailCollection := client.Database("emails").Collection("pending")
	webhookCollection := client.Database("webhooks").Collection("pending")

	defer func(){
		if err = client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	app := fiber.New()

	app.Use(logger.New())


	app.Post("/events/emails", func(c *fiber.Ctx) error {
		log.Println("Incoming")

		email := &EmailEvent{}

		c.BodyParser(email)

		fmt.Println(email)

		email.CreatedAt = time.Now()

		resp, err := emailCollection.InsertOne(context.Background(), email); if err != nil {

			c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"error": err,
			})

			return err
			
		}

		c.Status(http.StatusCreated).JSON(&fiber.Map{
			"data": resp,
		})

		return nil
	})


	app.Post("/events/webhooks", func(c *fiber.Ctx) error {
		webhook := &WebhookEvent{}

		c.BodyParser(webhook)

		webhook.CreatedAt = time.Now()

		resp, err := webhookCollection.InsertOne(context.Background(), webhook); if err != nil {
			c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"error": err,
			})

			return err;
		}

		c.Status(http.StatusAccepted).JSON(&fiber.Map{
			"data": resp,
		})

		return nil;
	})


	app.Listen(":9090")

}