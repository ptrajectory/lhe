package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Sender struct {
	Name string 
}

type Receiver struct {
	Email string
	URL string
}


type Event struct {
	Object interface{} 
	Type string 
	sender Sender
	receiver Receiver
}


func MainHandler(client *mongo.Client){
	log.Printf("------------------------------HANDLING NEW EVENTS----------------------------")
	// getting some data
	coll := client.Database("events").Collection("pending")

	fiveMinutesAgo := time.Now().Add(-5 *time.Minute)

	// last 5 minuts
	filter := bson.M{
		"created_at": bson.M{
			"$gte": fiveMinutesAgo.Format(time.RFC3339),
		},
	}

	cursor, err := coll.Find(context.Background(), filter)

	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.Background())

	var events []Event


	cursor.All(context.Background(), &events); if err == nil {
		fmt.Print("Something went wrong while serializing the data")
	}


	for i := 0; i < len(events); i++ {

		switch(events[i].Type){
			case "send_email":
				emailSender(events[i])
			case "send_webhook":
				webhookSender(events[i])
			default: 
				fmt.Println("There is nothing to do here.")		
		}

	}

	time.Sleep(time.Minute * 5);

	MainHandler(client)

}