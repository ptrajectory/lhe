package main

import (
	"context"
	"log"
	"os"
	handler "worker/handlers"

	"github.com/joho/godotenv"
	"github.com/resendlabs/resend-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



func main(){

	APP_ENV := os.Getenv("APP_ENV"); if APP_ENV == "DEVELOPMENT" || len(APP_ENV) == 0 {
		if err := godotenv.Load(); err != nil {
			panic("UNABLE TO LOAD ENV")
		}
	}

	MONGO_CONNECTION_STRING := os.Getenv("MONGO_CONNECTION_STRING"); if len(MONGO_CONNECTION_STRING) == 0 {
		panic("MONGO_CONNECTION_STRING NOT SET!!!")
	}
	
	RESEND_API_KEY := os.Getenv("RESEND_API_KEY"); if len(RESEND_API_KEY) == 0 {
		panic("RESEND_API_KEY NOT SET!!!")
	}



	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().ApplyURI(MONGO_CONNECTION_STRING).SetServerAPIOptions(serverAPI)


	client, err := mongo.Connect(context.TODO(), opts)

	if(err != nil) {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	resendClient := resend.NewClient(RESEND_API_KEY)

	log.Println(":::::::::::::::::::::::::::::::::::::::::::INITIALIZATION COMPLETE:::::::::::::::::::::::::::::::::::");

	handler.MainHandler(client, resendClient)

}