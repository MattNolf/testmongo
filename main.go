// Test the usage of decoding mongo bson documents, into objects
package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
)

const (
	mongoURI        = "mongodb://localhost:27017"
	mongoDB         = "mytest"
	mongoCollection = "mycollection"

	testCount = 1000
)

var (
	expectedFirstLine  = "expected_first_line"
	expectedSecondLine = "expected_second_line"

	failCount = 0
)

type Address struct {
	FirstLine  string `bson:"first_line"`
	SecondLine string `bson:"second_line"`
}

type User struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Name    string             `bson:"name"`
	Address Address            `bson:"address"`
	Age     int                `bson:"age"`
}

func main() {
	mongoClient, err := mongo.NewClient(mongoURI)
	if err != nil {
		return
	}

	err = mongoClient.Connect(context.Background())
	if err != nil {
		return
	}

	for i := 0; i < testCount; i++ {
		addr, err := GetAddress(mongoClient)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				failCount++
			}
		} else {
			fmt.Println(addr)
		}
	}
	fmt.Println(failCount)
}

func GetAddress(mongoClient *mongo.Client) (*User, error) {
	var address User

	singleResult := mongoClient.
		Database(mongoDB).
		Collection(mongoCollection).
		FindOne(
			context.Background(),
			bson.M{
				"address": bson.M{
					"first_line":  expectedFirstLine,
					"second_line": expectedSecondLine},
			})

	if singleResult.Err() != nil {
		return nil, singleResult.Err()
	}

	err := singleResult.Decode(&address)
	return &address, err
}
