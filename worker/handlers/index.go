package handler

import (
	"time"

	"github.com/resendlabs/resend-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type Sender struct {
	Name string 
	Id string
	Email string
}

type Receiver struct {
	Email string
	URL string
}

type WebhookReceiver struct {
	URL string
	// TODO: good place for any security headers
}

type EmailReceiver struct {
	Email string
}

type EmailPayload struct {
	Body string
	Subject string
	From string
}


type WebhookPayload struct {
	Type string
	Resource string
	Object interface{}
}

type WebhookEvent struct {
	Type string
	Sender Sender
	Receiver  WebhookReceiver
	CreatedAt time.Time `bson:"created_at"`
	Payload WebhookPayload
	ID string `bson:"_id"`
}

type EmailEvent struct {
	Sender Sender
	Receiver EmailReceiver
	Payload EmailPayload
	CreatedAt time.Time `bson:"created_at"`
	ID string `bson:"_id"`
	SendId string `bsim:"send_id" json:"send_id"`
}

type AppClients struct {
	client *mongo.Client
	resend *resend.Client
}	


func MainHandler(mongo *mongo.Client, resend *resend.Client){

	app := &AppClients{
		client: mongo,
		resend: resend,
	}

	WebhookHandler(app)
}