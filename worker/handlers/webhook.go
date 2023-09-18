package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func webhookSender(event WebhookEvent, c *AppClients)  {
	fmt.Println("::::::::::::::::::::::::WEBHOOKS::::::::::::::::::::::::::::")

	endpoint_url := event.Receiver.URL


	fmt.Printf("Payload:: %v \n\n", event.Payload)

	payloadBytes, err := json.Marshal(event.Payload); if err != nil {
		fmt.Printf("Unable to send request %s", err.Error())
		// return err
	}

	httpRequest, err := http.NewRequest(http.MethodPost,endpoint_url, bytes.NewBuffer(payloadBytes)); if err != nil {
		fmt.Printf("Unable to send request %s", err.Error())
		// return err
	}

	httpRequest.Header.Add("Content-Type", "application/json")

	client := http.Client{}

	resp, err := client.Do(httpRequest); if err != nil {
		fmt.Printf("Unable to send request %s", err.Error())
		// return err;
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body); if err != nil {
		fmt.Printf("Error reading response %s", err.Error())
		// return err;
	}

	fmt.Println("Response:", string(body))

	fmt.Printf("THE EVENT %s  %s\n\n", event.ID, event.Type)

	id, _ := primitive.ObjectIDFromHex(event.ID)

	filter := bson.D{{Key: "_id", Value: id}}

	deleted, err := c.client.Database("webhooks").Collection("pending").DeleteOne(context.Background(), filter); if err != nil {
		fmt.Printf("Something went wrong %v", err.Error())
		// return err;
	}
	fmt.Printf("The deleted %v \n\n", (*deleted).DeletedCount)

	if resp.StatusCode == 200 {

		_, err := c.client.Database("webhooks").Collection("success").InsertOne(context.Background(), event); if err != nil {
			fmt.Printf("Something went wrong %v", err.Error())
			// return err;
		}

	}else {
		_, err := c.client.Database("webhooks").Collection("failed").InsertOne(context.Background(), event); if err != nil {
			fmt.Printf("Something went wrong %v", err.Error())
			// return err
		}
 	}


	// return nil 
}


func WebhookHandler(app *AppClients) {

	pendingWebhooks := app.client.Database("webhooks").Collection("pending")


	filter := bson.M{
		"created_at": bson.M{
			"$gte": time.Now().Add(- time.Minute * 5),
		},
	}

	cursor, err := pendingWebhooks.Find(context.Background(), filter); if err != nil {
		fmt.Printf("Something went wrong %v", err.Error())
		// return err
	}

	webhooks := []WebhookEvent{}


	err = cursor.All(context.Background(), &webhooks); if err != nil {
		fmt.Printf("Something went wrong %v \n", err.Error())
		// return err
	}

	if(webhooks) == nil {
		emailHandler(app)
		return
	}

	for i :=0; i < len(webhooks); i++ {
		webhookSender((webhooks)[i], app)
	}	

	time.Sleep(time.Second * 30)
	emailHandler(app)

	// return nil;

}