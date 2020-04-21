package main

//Some golang handling of mongoDB database for
//restaurant DB in https://docs.mongodb.com/manual/tutorial/geospatial-tutorial/

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Resto struct {
	ID       primitive.ObjectID `"json":"_id" "bson":"_id,omitempty"`
	Location struct {
		Coordinates [2]float64 `"json":"coordinates" "bson":"coordinates"`
		Kind        string     `"json":"type" "bson":"type"`
	} `"json":"location" "bson":"location"`
	Name string `"json":"name" "bson":"name"`
}

func main() {
	const mongodbURI = "mongodb://localhost:27017/"
	var ClientOptions = options.Client().ApplyURI(mongodbURI)
	/*
		var test1 = bson.M{
			"location" : {
				"coordinates" : [5,5],
				"type" : "point"
			}
			"name" : "test"
		}
	*/
	client, err := mongo.Connect(context.TODO(), ClientOptions)
	if err != nil {
		log.Fatal("cannot connect to mongodb")
	}
	defer client.Disconnect(context.TODO())

	//neighborhoods := client.Database("test").Collection("neighborhoods")
	filter := bson.D{{}}
	restaurants := client.Database("test").Collection("restaurants")

	//	res := restaurants.FindOne(context.TODO(), filter)

	cursor, err := restaurants.Find(context.TODO(), filter, options.Find().SetLimit(5))

	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cursor.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Resto
		err := cursor.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%+v", elem)

	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

}
