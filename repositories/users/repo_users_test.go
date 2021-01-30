package users

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/unnamedxaer/gymm-api/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ur *UserRepository
	db *mongo.Database
	u  userData
)

func TestMain(m *testing.M) {
	logger := zerolog.New(os.Stdout)
	repositories.EnsureTestEnv()

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		panic("environment variable 'DB_NAME' is not set")
	}
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		panic("environment variable 'MONGO_URI' is not set")
	}
	db, err := repositories.GetDatabase(&logger, mongoURI, dbName)
	if err != nil {
		panic(err)
	}

	err = repositories.CreateCollections(&logger, db)
	if err != nil {
		panic(err)
	}
	defer repositories.DisconnectDB(&logger, db)

	usersCol := db.Collection("users")
	_, err = usersCol.DeleteMany(nil, nil)
	if err != nil {
		panic(err)
	}

	ur = NewRepository(&zerolog.Logger{}, usersCol)

	password, _ := hashPassword("TheSecretestPasswordEver123$%^")
	u = userData{
		Username:     "John Silver",
		EmailAddress: "johnsilver@email.com",
		Password:     password,
		CreatedAt:    time.Now().UTC(),
	}

	code := m.Run()
	os.Exit(code)
}

func TestCreateUser(t *testing.T) {
	gotUser, err := ur.CreateUser(u.Username, u.EmailAddress, u.Password)
	if err != nil {
		t.Fatal(err)
	}
	if u.EmailAddress != gotUser.EmailAddress ||
		u.Username != gotUser.Username {
		t.Errorf("Expect to get user base on data: %v, got: %v",
			u,
			gotUser)
	}

	if gotUser.ID == "" {
		t.Fatalf("Expected 'ID' to not be empty, got '%s' for email: %s",
			gotUser.ID,
			u.EmailAddress)
	}

	if gotUser.CreatedAt.IsZero() {
		t.Fatalf("Expected 'CreateAt' to not be zero value, got '%s' for email: %s",
			gotUser.CreatedAt,
			u.EmailAddress)
	}
}

func TestCreateUserDuplicatedEmail(t *testing.T) {
	gotUser, err := ur.CreateUser(u.Username, u.EmailAddress, u.Password)
	if err != nil {
		if errors.Is(err, repositories.NewErrorEmailAddressInUse()) {
			return
		}
		t.Fatal(err)
	}
	t.Fatalf("Expected to get error '%s' for email addres: '%s', got new user with ID: '%s'",
		repositories.NewErrorEmailAddressInUse(),
		u.EmailAddress,
		gotUser.ID)
}

func TestGetUserByID(t *testing.T) {

	results, err := ur.col.InsertOne(context.TODO(), u)
	if err != nil {
		t.Fatal(err)
	}

	uID := results.InsertedID.(primitive.ObjectID).Hex()

	gotUser, err := ur.GetUserByID(uID)
	// gotUser, err := ur.GetUserByID("6010835ea16aa7173e36fc0d")
	if err != nil {
		t.Fatal(err)
	}

	if timesEqual(u.CreatedAt, gotUser.CreatedAt) == false ||
		u.EmailAddress != gotUser.EmailAddress ||
		u.Username != gotUser.Username ||
		uID != gotUser.ID {
		t.Errorf("Expect to get user like: %v, got: %v", u, gotUser)
	}
}

func TestGetUserByIDNotExisting(t *testing.T) {
	uID := "60108393da81e60598d5347f"
	uObjectID, _ := primitive.ObjectIDFromHex(uID)
	ur.col.DeleteOne(context.TODO(), bson.M{"_id": uObjectID})

	gotUser, err := ur.GetUserByID(uID)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return
		}
		t.Fatal(err)
	}
	t.Errorf("Expect to NOT get any user for _id: '%s', but got: %v", uID, gotUser)
}

func timesEqual(t1, t2 time.Time) bool {

	return t1.Year() == t2.Year() &&
		t1.Month() == t2.Month() &&
		t1.Day() == t2.Day() &&
		t1.Hour() == t2.Hour() &&
		t1.Minute() == t2.Minute() &&
		t1.Second() == t2.Second()
}
