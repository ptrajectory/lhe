package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/resendlabs/resend-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)



func emailSender(event EmailEvent, c *AppClients) {


	params := resend.SendEmailRequest{
		From: event.Sender.Email,
		To: []string {event.Receiver.Email},
		Html: event.Payload.Body,
		Subject: event.Payload.Subject,
	}

	sent, err := c.resend.Emails.Send(&params); if err != nil {
		fmt.Printf("Something went wrong :: %v", err.Error())
		event.SendId = sent.Id

		c.client.Database("emails").Collection("failed").InsertOne(context.TODO(), event);
	}

	event.SendId = sent.Id
	c.client.Database("emails").Collection("success").InsertOne(context.Background(), event)


	fmt.Println(sent.Id)

	id, _ := primitive.ObjectIDFromHex(event.ID)

	delete_filter := bson.D{{"_id", id}}

	deleted, err := c.client.Database("emails").Collection("pending").DeleteOne(context.TODO(), delete_filter); if err != nil {
		fmt.Printf("Something went wrong:: %s", err)
	}

	fmt.Printf("Deleted %d", deleted.DeletedCount)

}

func emailHandler(app *AppClients) {
	fmt.Println("::::::::::::::::::::::::EMAILS::::::::::::::::::::::::::::")
	pendingEmails := app.client.Database("emails").Collection("pending")

	filter := bson.M{
		"created_at": bson.M{
			"$gte": time.Now().Add(-time.Minute * 5),
		},
	}


	cursor, err := pendingEmails.Find(context.Background(), filter); if err != nil {
		fmt.Printf("Something went wrong %v", err)
		// return err;
	}

	emails := []EmailEvent{}

	err = cursor.All(context.Background(), &emails); if err != nil{
		fmt.Printf("Something went wrong %v", err)
		// return err;
	}

	fmt.Printf("HERE Are the emails %v \n", emails)

	if((emails) == nil){
		WebhookHandler(app)
		return
	}

	for i := 0; i < len(emails); i++ {
		emailSender((emails)[i], app)
	}

	time.Sleep(time.Second * 30)
	WebhookHandler(app)

	// return nil
}