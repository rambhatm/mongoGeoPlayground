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

type Restaurant struct {
	ID       primitive.ObjectID `"bson":"_id,omitempty"`
	Location struct {
		Coordinates [2]float64 `"bson":"coordinates"`
		Kind        string     `"bson":"type"`
	} `"bson":"location"`
	Name string `"bson":"name"`
}

func createGeoIndex(coll *mongo.Collection) {
	model := mongo.IndexModel{
		Keys: bson.M{
			"location": "2dsphere",
		}, Options: nil,
	}
	coll.Indexes().CreateOne(context.TODO(), model)
}

func main() {
	const mongodbURI = "mongodb://localhost:27017/"
	var ClientOptions = options.Client().ApplyURI(mongodbURI)

	client, err := mongo.Connect(context.TODO(), ClientOptions)
	if err != nil {
		log.Fatal("cannot connect to mongodb")
	}
	defer client.Disconnect(context.TODO())

	filter := bson.D{{}}
	restaurants := client.Database("test").Collection("restaurants")

	createGeoIndex(restaurants)

	cursor, err := restaurants.Find(context.TODO(), filter, options.Find().SetLimit(5))

	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cursor.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Restaurant
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
