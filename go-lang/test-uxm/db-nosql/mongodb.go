
// GO Lang :: Sample MongoDB

package main

import (
	"context"
	"time"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)


type Person struct {
	Name string
	Age  int
	City string
}



func main() {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	// To do this in a single step, you can use the Connect function:
//	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
//	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	// Calling Connect does not block for server discovery.
	// If you wish to know if a MongoDB server has been found and connected to, use the Ping method:
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalln(err)
	}

	collection := client.Database("smart_framework").Collection("goTest")

	// INSERT: The Collection instance can then be used to insert documents:
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	res, err := collection.InsertOne(ctx, bson.M{"name": "pi", "value": 3.14159, "views": 0})
	id := res.InsertedID

	fmt.Println("InsertedID: ", id)


	// Raw string representation of the MongoDB doc _id
//	idStr := "5d2399ef96fb765873a24bae"
	// Create a BSON ObjectID by passing string to ObjectIDFromHex() method
//	docID, err := primitive.ObjectIDFromHex(idStr)

	// FIND MANY: Several query methods return a cursor, which can be used like this:
	theQuery := bson.M{"name": bson.M{ "$in": bson.A{"pi", "qr"} } }
	ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := collection.Find(ctx, theQuery)
	if err != nil { log.Fatal(err) }
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil { log.Fatal(err) }
		// do something with result....
		fmt.Printf("Found document: %+v\n", result)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}


	// FIND ONE:
	var result struct {
		ID    primitive.ObjectID `bson:"_id" json:"id,omitempty"`
		Value float64
	}
	filter := bson.M{"name": "pi"}
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	// Do something with result...
	fmt.Printf("Found a single document: %+v\n", result)
	fmt.Println("Document.Value: ", result.Value)
	fmt.Println("Document.ID as Object: ", result.ID)
	fmt.Println("Document.ID as String: ", result.ID.Hex())

	// UPDATE:
	filter2 := bson.D{{"name", "pi"}}
	update := bson.D{
		{"$inc", bson.D{
			{"views", 1},
		}},
	}
	updateResult, err := collection.UpdateOne(context.TODO(), filter2, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	// DELETE
	delete := bson.D{{"views", 0}}
	deleteResult, err := collection.DeleteMany(context.TODO(), delete)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)

	// INSERT MANY
	ruan := Person{"Ruan", 34, "Cape Town"}
	james := Person{"James", 32, "Nairobi"}
	frankie := Person{"Frankie", 31, "Nairobi"}
	trainers := []interface{}{james, frankie, ruan}
	insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)

	// UPDATE
	filter3 := bson.D{{}}
	update2 := bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
	}
	updateResult2, err := collection.UpdateOne(context.TODO(), filter3, update2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult2.MatchedCount, updateResult2.ModifiedCount)


}

// END
