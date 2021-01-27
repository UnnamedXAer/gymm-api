package repositories

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	usersCollectionName     = "users"
	trainingsCollectionName = "trainings"
)

// GetDatabase connects to mongodb and returns database
func GetDatabase(l *zerolog.Logger, uri, dbName string) (*mongo.Database, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, errors.WithMessagef(err, "connect to db '%s'", dbName)
	}
	l.Info().Msgf("connected to db '%s'", dbName)

	db := client.Database(dbName)
	return db, nil
}

// CreateCollections creates mongodb collections
func CreateCollections(l *zerolog.Logger, db *mongo.Database) error {
	collections, err := db.ListCollectionNames(nil, nil)
	if err != nil {
		return err
	}
	if sliceIndexOf(collections, usersCollectionName) == -1 {
		err = createUsersCollection(l, db, usersCollectionName)
		if err != nil {
			return err
		}
	} else {
		l.Info().Msgf("collection '%s' already exists - skipped")
	}
	if sliceIndexOf(collections, trainingsCollectionName) == -1 {
		err = createTrainingCollection(l, db, trainingsCollectionName)
		return err
	} else {
		l.Info().Msgf("collection '%s' already exists - skipped")
	}
	return nil
}

func createUsersCollection(l *zerolog.Logger, db *mongo.Database, collectionName string) error {
	err := db.CreateCollection(context.Background(), collectionName, &options.CreateCollectionOptions{
		Collation: &options.Collation{
			Strength: 2,
		},
	})
	if err != nil {
		return errors.WithMessagef(err, "create '%s' collection", collectionName)
	}
	l.Info().Msgf("collection '%s' created", collectionName)

	col := db.Collection(collectionName)

	emailIndexName := "unique email_address"
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email_address", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("unique email_address"),
	}

	indexName, err := col.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return errors.WithMessagef(err, "create index '%s' on '%s' collection", emailIndexName, collectionName)
	}

	l.Info().Msgf("index '%s' on collection '%s' created", indexName, collectionName)
	return nil
}

func createTrainingCollection(l *zerolog.Logger, db *mongo.Database, collectionName string) error {
	const collectionName string = "trainings"
	err := db.CreateCollection(context.Background(), collectionName)
	if err != nil {
		l.Info().Msgf("collection '%s' created", collectionName)
	}
	l.Info().Msgf("collection '%s' created", collectionName)
	return nil
}

func sliceIndexOf(slice []string, s string) int {
	for i, count := 0, len(slice); i < count; i++ {
		if slice[i] == s {
			return i
		}
	}
	return -1
}
