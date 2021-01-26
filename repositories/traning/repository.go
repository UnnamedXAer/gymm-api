package traning

import (
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

type TrainingRepository struct {
	col *mongo.Collection
	l   zerolog.Logger
}

func NewRepository(logger zerolog.Logger, collection *mongo.Collection) *TrainingRepository {
	return &TrainingRepository{
		collection,
		logger,
	}
}

// func makeme(uri string) error {
// 	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
// 	if err != nil {
// 		return err
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	err = client.Connect(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	dbName, ok := os.LookupEnv("DB_NAME")
// 	if ok == false {
// 		log.Println("db name not specified fallback to default 'gymm-api'")
// 		dbName = "gymm-api"
// 	}
// 	log.Println("DB name: ", dbName)
// 	db := client.Database(dbName)
// 	coll := db.Collection("trainings")
// 	var r *TrainingRepository = NewRepository(zerolog.New(os.Stdout), coll)

// 	return nil
// }
