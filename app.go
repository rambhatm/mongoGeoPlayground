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
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Location struct {
		Coordinates [2]float64 `bson:"coordinates"`
		Kind        string     `bson:"type"`
	} `bson:"location"`
	Name string `bson:"name"`
}

//https://docs.mongodb.com/manual/reference/method/db.collection.createIndex/#recreating-an-existing-index
func createGeoIndex(coll *mongo.Collection) {
	model := mongo.IndexModel{
		Keys: bson.M{
			"location": "2dsphere",
		}, Options: nil,
	}
	coll.Indexes().CreateOne(context.TODO(), model)
}

func searchWithinRadius(coll *mongo.Collection, coords [2]float64, radiusInM int) {
	log.Printf("\nSearching [%f,%f] radius %d", coords[0], coords[1], radiusInM)

	nearSphere := bson.M{"$geometry": bson.M{"coordinates": coords, "type": "Point"}, "$maxDistance": radiusInM}
	filter := bson.M{"location": bson.M{"$nearSphere": nearSphere}}

	cursor, err := coll.Find(context.TODO(), filter, options.Find().SetLimit(5))

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

func main() {
	const mongodbURI = "mongodb://localhost:27017/"
	var ClientOptions = options.Client().ApplyURI(mongodbURI)

	client, err := mongo.Connect(context.TODO(), ClientOptions)
	if err != nil {
		log.Fatal("cannot connect to mongodb")
	}
	defer client.Disconnect(context.TODO())

	restaurants := client.Database("test").Collection("restaurants")

	createGeoIndex(restaurants)
	testLoc := [2]float64{-73.871194, 40.6730975}
	searchWithinRadius(restaurants, testLoc, 50)
	searchWithinRadius(restaurants, testLoc, 100)
	searchWithinRadius(restaurants, testLoc, 200)
	searchWithinRadius(restaurants, testLoc, 500)
	searchWithinRadius(restaurants, testLoc, 1000)

}
