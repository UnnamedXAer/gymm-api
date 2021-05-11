package repositories

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/testhelpers"
	"github.com/unnamedxaer/gymm-api/usecases"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	db         *mongo.Database
	loggerMock zerolog.Logger
)

const (
	colSuffix = "_test_create_col"
)

func TestMain(m *testing.M) {
	testhelpers.EnsureTestEnv()
	loggerMock = zerolog.New(nil)

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatalln("environment variable 'DB_NAME' is not set")
	}
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatalln("environment variable 'MONGO_URI' is not set")
	}
	var err error
	db, err = GetDatabase(&loggerMock, mongoURI, dbName)
	if err != nil {
		log.Fatalln(err)
	}

	filter := bson.D{}
	collections, err := db.ListCollectionNames(context.Background(), filter)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(collections)
	ctx := context.Background()
	for _, col := range collections {
		if !strings.HasSuffix(col, colSuffix) {
			continue
		}
		err = db.Collection(col).Drop(ctx)
		if err != nil {
			log.Fatalln(err)
		}
	}

	os.Exit(m.Run())
}

func TestCreateUsersCollection(t *testing.T) {
	colName := UsersCollectionName + colSuffix
	err := createUsersCollection(&loggerMock, db, colName)
	if err != nil {
		t.Fatal(err)
	}
	usCol := db.Collection(colName)
	n := 0
	getN := func() int {
		n++
		return n
	}

	tmplName := "username >> %d <<"
	input := bson.M{
		"username":      fmt.Sprintf(tmplName, getN()),
		"email_address": "gopher@example.com",
		"password":      "not_hashed_password",
		"created_at":    time.Now(),
	}
	_, err = usCol.InsertOne(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}

	inputUnique := input
	inputUnique["username"] = fmt.Sprintf(tmplName, getN())
	inputUnique["email_address"] = "mongopher@example.com"
	_, err = usCol.InsertOne(context.Background(), inputUnique)
	if err != nil {
		t.Fatalf("want error to be %v, got %v", nil, err)
	}

	inputDuplicate := input
	inputDuplicate["username"] = fmt.Sprintf(tmplName, getN())
	_, err = usCol.InsertOne(context.Background(), input)
	if !usecases.IsDuplicatedError(err) {
		t.Fatalf("want error like: %s, got %v", "E11000 duplicate key error collection", err)
	}

	result, err := usCol.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		t.Fatal(err)
	}

	wantCnt := int64(n - 1)
	if result != wantCnt {
		t.Fatalf("want %d documents in collection, got %d", wantCnt, result)
	}
}

func TestCreateExercisesCollection(t *testing.T) {
	colName := ExercisesCollectionName + colSuffix
	err := createExercisesCollection(&loggerMock, db, colName, false)
	if err != nil {
		t.Fatal(err)
	}
	exCol := db.Collection(colName)
	n := 0
	getN := func() int {
		n++
		return n
	}
	tmplDesc := "Deadlift is compound movement. >> %d <<"
	input := bson.M{
		"name":        "Deadlift",
		"description": fmt.Sprintf(tmplDesc, getN()),
		"set_unit":    1,
		"created_at":  time.Now(),
		"created_by":  primitive.ObjectID([12]byte{})}

	_, err = exCol.InsertOne(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}

	inputUnique := input
	inputUnique["name"] = "OHP"
	inputUnique["description"] = fmt.Sprintf(tmplDesc, getN())
	_, err = exCol.InsertOne(context.Background(), inputUnique)
	if err != nil {
		t.Fatalf("want error to be %v, got %v", nil, err)
	}

	inputDuplicate := input
	inputDuplicate["description"] = fmt.Sprintf(tmplDesc, getN())
	_, err = exCol.InsertOne(context.Background(), input)
	if err == nil {
		t.Fatalf("want index violation error, got %v", err)
	}

	if !usecases.IsDuplicatedError(err) {
		t.Fatalf("want error like: %s, got %v", "E11000 duplicate key error collection", err)
	}

	result, err := exCol.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		t.Fatal(err)
	}

	wantCnt := int64(n - 1)
	if result != wantCnt {
		t.Fatalf("want %d documents in collection, got %d", wantCnt, result)
	}

	filter := bson.M{"$text": bson.M{"$search": "OHP"}}
	_, err = exCol.Find(context.Background(), filter)
	if err != nil {
		if err.Error() == "text index required for $text query" {
			t.Fatalf("want text index to exists on exercise collection")
		}
		t.Fatal(err)
	}
}

func TestCreateTrainingsCollection(t *testing.T) {
	colName := TrainingsCollectionName + colSuffix
	err := createTrainingsCollection(&loggerMock, db, colName)
	if err != nil {
		t.Fatal(err)
	}
	trCol := db.Collection(colName)

	input := bson.M{
		"user_id":    primitive.ObjectID([12]byte{}),
		"start_time": time.Now().Add(-1 + time.Hour),
		"end_time":   time.Now(),
		"comment":    "too heavy",
		"exercises":  bson.A{},
	}

	_, err = trCol.InsertOne(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}

	result, err := trCol.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		t.Fatal(err)
	}

	var wantCnt int64 = 1
	if result != wantCnt {
		t.Fatalf("want %d documents in collection, got %d", wantCnt, result)
	}
}

func TestCreateResPwdReqCollection(t *testing.T) {
	ctx := context.TODO()
	colName := ResPwdReqCollectionName + colSuffix
	err := createResPwdReqCollection(&loggerMock, db, colName)
	if err != nil {
		t.Fatal(err)
	}
	rprCol := db.Collection(colName)

	input := bson.M{
		"emailAddress": "",
		"expires_at":   time.Now(),
		"status":       entities.ResetPwdStatusNoActionYet,
		"comment":      "",
		"created_at":   time.Now().Add(-15 * time.Minute),
	}

	_, err = rprCol.InsertOne(ctx, input)
	if err != nil {
		t.Fatal(err)
	}

	result, err := rprCol.CountDocuments(ctx, bson.M{})
	if err != nil {
		t.Fatal(err)
	}

	var wantCnt int64 = 1
	if result != wantCnt {
		t.Fatalf("want %d documents in collection, got %d", wantCnt, result)
	}
}
