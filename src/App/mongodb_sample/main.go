package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"tutApp/mongodb_driver_tut/mymongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gin-gonic/gin"
	// "net/http"
)

type Trainer struct {
	Name string
	Age  int
	City string
	Tel  string
}

type Car struct {
	Brand string
	Price int
}

func main() {
	// create()
	// delete()
	// findMany()
	// find()
	// update()


	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/find", func(c *gin.Context) {
		result := find()
		fmt.Println(result)
		c.JSON(200, gin.H{
			"message": result,
		})
	})
	r.Run(":9999") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func update() {

	client, err := mymongo.Init()
	collection_all_redeem_log := client.Database("grid3x3").Collection("all_redeem_log")

	if err != nil {
		log.Fatal(err)
	}

	filter := bson.D{{"name", "Ana cistla"}}
	update := bson.D{
		{"$set", bson.D{{"age", 37}, {"tel", "0999"}}},
	}

	// updateResult, err := collection_all_redeem_log.ReplaceOne(
	// 	context.TODO(),
	// 	bson.M{"author": "Nicolas Raboy"},
	// 	bson.M{
	// 		"name":  "Ana cistla",
	// 		"author": "Nicolas Raboy",
	// 	},
	// )

	updateResult, err := collection_all_redeem_log.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(updateResult)
}

func find() string {
	// create a value into which the result can be decoded
	var result Trainer
	client, err := mymongo.Init()
	collection_all_redeem_log := client.Database("grid3x3").Collection("all_redeem_log")

	err = collection_all_redeem_log.FindOne(context.TODO(), bson.D{{}}).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	json_data ,err := json.Marshal(result)

	fmt.Printf("Found a single document: %+v\n", result)
	return string(json_data);
}

func findMany() {

	client, err := mymongo.Init()
	collection_all_redeem_log := client.Database("grid3x3").Collection("all_redeem_log")

	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetLimit(4)

	// Here's an array in which you can store the decoded documents
	var results []*Trainer

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection_all_redeem_log.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Trainer
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)

	for key, value := range results {
		fmt.Println(key, *value)
	}

}

func delete() {
	client, err := mymongo.Init()
	// fmt.Println("Connected to MongoDB!")

	collection_all_redeem_log := client.Database("grid3x3").Collection("all_redeem_log")
	deleteResult, err := collection_all_redeem_log.DeleteMany(context.TODO(), bson.D{{"name", "Ana"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
}

func create() {
	client, err := mymongo.Init()
	// fmt.Println("Connected to MongoDB!")

	collection_all_redeem_log := client.Database("grid3x3").Collection("all_redeem_log")

	// ash := Trainer{"Ash", 10, "Pallet Town"}

	benz := Car{Brand: "Benz", Price: 3000000}

	misty := Trainer{Name: "Ana", Age: 45, City: "Moscow City"}
	brock := Trainer{Name: "Max", Age: 33, City: "St.Petersburg City"}
	trainers := []interface{}{misty, brock, benz}

	insertManyResult, err := collection_all_redeem_log.InsertMany(context.TODO(), trainers)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
}
