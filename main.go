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

	user := User{
		Name: "some-name",
		Address: Address{
			FirstLine:  expectedFirstLine,
			SecondLine: expectedSecondLine,
		},
		Age: 1,
	}

	_, err = AddUser(mongoClient, &user)
	if err != nil {
		return
	}

	for i := 0; i < testCount; i++ {
		_, err := GetUser(mongoClient)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				failCount++
			}
		}
	}
	fmt.Println(failCount)
	err = Cleanup(mongoClient)
}

func GetUser(mongoClient *mongo.Client) (*User, error) {
	var address User
	//filter := bson.M{
	//	"address": bson.M{
	//		"first_line":  expectedFirstLine,
	//		"second_line": expectedSecondLine}}

	filter := bson.M{
		"address.first_line":  expectedFirstLine,
		"address.second_line": expectedSecondLine,
	}

	singleResult := mongoClient.
		Database(mongoDB).
		Collection(mongoCollection).
		FindOne(
			context.Background(), filter)

	if singleResult.Err() != nil {
		return nil, singleResult.Err()
	}

	err := singleResult.Decode(&address)
	return &address, err
}

func AddUser(mongoClient *mongo.Client, user *User) (*mongo.InsertOneResult, error) {
	return mongoClient.
		Database(mongoDB).
		Collection(mongoCollection).
		InsertOne(context.Background(), user)
}

func Cleanup(mongoClient *mongo.Client) error {
	filter := bson.M{}

	_, err := mongoClient.
		Database(mongoDB).
		Collection(mongoCollection).
		DeleteMany(context.Background(), filter)

	return err
}
