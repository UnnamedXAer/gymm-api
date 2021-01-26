package repository

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/unnamedxaer/gymm-api/server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func (rep *MongoRepository) Initialize(uri string) error {
	var err error
	rep.Client, err = mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return err

	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = rep.Client.Connect(ctx)
	if err != nil {
		return err
	}

	dbName, ok := os.LookupEnv("DB_NAME")
	if ok == false {
		log.Println("db name not specified fallback to default 'gymm-api'")
		dbName = "gymm-api"
	}
	log.Println("DB name: ", dbName)
	rep.DB = rep.Client.Database(dbName)

	return nil
}

func (rep *MongoRepository) CreateUser(u *models.User) error {
	col := rep.DB.Collection("users")
	results, err := col.InsertOne(context.TODO(), u)
	if err != nil {
		return err
	}
	u.ID = results.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (rep *MongoRepository) GetUserById(u *models.User) error {
	log.Println("[mongo - GetUserById] " + u.ID)
	oID, err := primitive.ObjectIDFromHex(u.ID)
	if err != nil {
		return err
	}
	col := rep.DB.Collection("users")
	err = col.FindOne(context.TODO(), bson.M{"_id": oID}).Decode(&u)
	return err
}
