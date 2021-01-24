package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/unnamedxaer/gymm-api/server"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var app = server.App{}

func TestMain(m *testing.M) {
	godotenv.Overload(".test.env")
	log.Println(" [TestMain] start " + os.Getenv("ENV"))

	app.InitializeApp()
	clearTables("mongo")

	code := m.Run()

	log.Println(" [TestMain] end")
	os.Exit(code)
}

func clearTables(repoType string) {
	log.Println(" [clearTables] start - " + repoType)
	switch repoType {
	case "mongo":
		clearTablesMongo()
	default:
		log.Fatalln("unsupported repo type")
	}
}

func clearTablesMongo() {
	log.Println(" [clearTablesMongo] start")
	db, err := getDb("mongo")
	mongoDB := db.(*mongo.Database)

	if err != nil {
		log.Fatal(err)
	}
	usersColl := mongoDB.Collection("users")
	// results, err := usersColl.DeleteMany(context.Background(), bson.M{})
	results, err := usersColl.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(" [clearTablesMongo] results " + fmt.Sprintf("%v", results))
}

func getDb(provider string) (interface{}, error) {

	switch provider {
	case "mongo":
		return getMongoDB()
	}
	return nil, fmt.Errorf("unsupported provider '%s'", provider)
}

func getMongoDB() (*mongo.Database, error) {
	log.Println("[getMongoDB] start")
	uri := os.Getenv("MONGO_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err

	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	dbName, ok := os.LookupEnv("DB_NAME")
	if ok == false {
		log.Panic("db name not specified")
	}
	log.Println("DB name: ", dbName)

	return client.Database(dbName), nil
}
