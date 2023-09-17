package main

import (
	"context"
	handler "lhe/handlers"
	
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Event struct {
	Object interface{} 
	Type string 
}



func main(){


	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().ApplyURI("mongodb+srv://lhe-admin:uX3yzkhUJhPQyjLT@lhe-test-cluster.dfpbft1.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)


	client, err := mongo.Connect(context.TODO(), opts)

	if(err != nil) {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	handler.MainHandler(client)

}