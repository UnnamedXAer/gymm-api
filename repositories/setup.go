package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/unnamedxaer/gymm-api/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	UsersCollectionName     = "users"
	TrainingsCollectionName = "trainings"
	ExercisesCollectionName = "exercises"
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

func GetCollection(l *zerolog.Logger, db *mongo.Database, collName string) *mongo.Collection {
	switch collName {
	// case exercisesCollectionName:
	// fallthrough
	case UsersCollectionName:
		fallthrough
	case TrainingsCollectionName:
		return db.Collection(collName)
	default:
		panic(fmt.Sprintf("unknown collection name '%s'", collName))
	}
}

// CreateCollections creates mongodb collections
func CreateCollections(l *zerolog.Logger, db *mongo.Database) error {
	collections, err := db.ListCollectionNames(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	colName := UsersCollectionName
	if helpers.StrSliceIndexOf(collections, colName) == -1 {
		err = createUsersCollection(l, db, colName)
		if err != nil {
			return err
		}
	} else {
		l.Info().Msgf("collection '%s' already exists - skipped", colName)
	}

	colName = TrainingsCollectionName
	if helpers.StrSliceIndexOf(collections, colName) == -1 {
		err = createTrainingsCollection(l, db, colName)
		return err
	} else {
		l.Info().Msgf("collection '%s' already exists - skipped", colName)
	}

	colName = ExercisesCollectionName
	if helpers.StrSliceIndexOf(collections, colName) == -1 {
		err = createExercisesCollection(l, db, colName)
		return err
	} else {
		l.Info().Msgf("collection '%s' already exists - skipped", colName)
	}
	return nil
}

func createUsersCollection(l *zerolog.Logger, db *mongo.Database, collectionName string) error {
	err := db.CreateCollection(context.Background(), collectionName, &options.CreateCollectionOptions{
		Collation: &options.Collation{
			Strength: 2,
			Locale:   "en",
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

func createTrainingsCollection(l *zerolog.Logger, db *mongo.Database, collectionName string) error {
	err := db.CreateCollection(context.Background(), collectionName)
	if err != nil {
		return errors.WithMessagef(err, "create '%s' collection", collectionName)
	}
	l.Info().Msgf("collection '%s' created", collectionName)
	return nil
}

func createExercisesCollection(l *zerolog.Logger, db *mongo.Database, collectionName string) error {
	err := db.CreateCollection(context.Background(), collectionName, &options.CreateCollectionOptions{
		Collation: &options.Collation{
			Strength: 2,
			Locale:   "en",
		},
	})
	if err != nil {
		return errors.WithMessagef(err, "create %q collection", collectionName)
	}
	l.Info().Msgf("collection %q created", collectionName)

	col := db.Collection(collectionName)

	emailIndexName := "name-set_unit"
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}, {Key: "time", Value: 1}},
		Options: options.Index().SetUnique(true).SetName(emailIndexName),
	}

	indexName, err := col.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return errors.WithMessagef(err, "create index %q on %q collection", emailIndexName, collectionName)
	}

	l.Info().Msgf("index %q on collection %q created", indexName, collectionName)
	return nil
}
