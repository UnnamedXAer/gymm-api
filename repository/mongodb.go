package repository

import (
	"context"
	"time"

	"github.com/unnamedxaer/gymm-api/server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	Ctx    context.Context
	Client *mongo.Client
	DB     *mongo.Database
}

func (rep *MongoRepository) Initialize(uri string) error {
	var err error
	rep.Client, err = mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return err

	}

	rep.Ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	if err != nil {
		return err
	}

	err = rep.Client.Connect(rep.Ctx)
	if err != nil {
		return err
	}

	rep.DB = rep.Client.Database("gymm-api")

	return nil
}

func (rep *MongoRepository) CreateUser(u *models.User) error {
	col := rep.DB.Collection("users")
	_, err := col.InsertOne(rep.Ctx, u)
	return err
}

func (rep *MongoRepository) GetUserById(u *models.User) error {
	col := rep.DB.Collection("users")
	err := col.FindOne(rep.Ctx, bson.M{"_id": u.ID}).Decode(&u)
	return err
}
