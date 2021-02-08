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
	_, err = usersCol.DeleteMany(nil, bson.D{})
	if err != nil {
		panic(err)
	}

	ur = NewRepository(&zerolog.Logger{}, usersCol)

	// password := []byte("TheSecretestPasswordEver123$%^")
	u = userData{
		Username:     "John Silver",
		EmailAddress: "johnsilver@email.com",
		Password:     []byte("TheSecretestPasswordEver123$%^"),
		CreatedAt:    time.Now().UTC(),
	}

	code := m.Run()
	os.Exit(code)
}

func TestCreateUser(t *testing.T) {
	clearCollection(t)
	gotUser, err := ur.CreateUser(u.Username, u.EmailAddress, u.Password)
	if err != nil {
		t.Fatal(err)
	}
	if u.EmailAddress != gotUser.EmailAddress ||
		u.Username != gotUser.Username {
		t.Fatalf("Expect to get user base on data: %v, got: %v",
			u,
			gotUser)
	}

	if gotUser.ID == "" {
		t.Fatalf("Expected 'ID' to NOT be zero value, got %q for email: %q",
			gotUser.ID,
			u.EmailAddress)
	}

	if gotUser.CreatedAt.IsZero() {
		t.Fatalf("Expected 'CreateAt' to NOT be zero value, got %q for email: %q",
			gotUser.CreatedAt,
			u.EmailAddress)
	}
}

func TestCreateUserDuplicatedEmail(t *testing.T) {
	clearCollection(t)
	insertMockUser(t)
	gotUser, err := ur.CreateUser(u.Username, u.EmailAddress, u.Password)
	if err != nil {
		if errors.Is(err, repositories.NewErrorEmailAddressInUse()) {
			return
		}
		t.Fatal(err)
	}
	t.Fatalf("Expected to get error %q for email addres: %q, got new user with ID: %q",
		repositories.NewErrorEmailAddressInUse(),
		u.EmailAddress,
		gotUser.ID)
}

func TestGetUserByID(t *testing.T) {
	clearCollection(t)
	results, err := ur.col.InsertOne(context.TODO(), u)
	if err != nil {
		t.Fatal(err)
	}

	uID := results.InsertedID.(primitive.ObjectID).Hex()

	gotUser, err := ur.GetUserByID(uID)
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
	clearCollection(t)
	uID := "60108393da81e60598d5347f"

	gotUser, err := ur.GetUserByID(uID)
	if err == nil || errors.Is(err, repositories.NewErrorNotFoundRecord()) == false {
		t.Fatalf("Expected to get error: %q, got: %v", repositories.NewErrorNotFoundRecord(), err)
	}
	if gotUser.ID != "" || gotUser.EmailAddress != "" {
		t.Errorf("Expect to NOT get any user for _id: %q, but got: %v", uID, gotUser)
	}
}

func TestGetUserByIDInvalidID(t *testing.T) {
	clearCollection(t)
	uID := "6s108393da81e60598d5347f"

	gotUser, err := ur.GetUserByID(uID)
	if err == nil || errors.Is(err, repositories.NewErrorInvalidID(uID)) == false {
		t.Fatalf("Expected to get error: %q, got: %v", repositories.NewErrorInvalidID(uID), err)
	}
	if gotUser.ID != "" || gotUser.EmailAddress != "" {
		t.Errorf("Expect to NOT get any user for _id: %q, but got: %v", uID, gotUser)
	}
}

func timesEqual(t1, t2 time.Time) bool {

	return t1.Year() == t2.Year() &&
		t1.Month() == t2.Month() &&
		t1.Day() == t2.Day() &&
		t1.Hour() == t2.Hour() &&
		t1.Minute() == t2.Minute() &&
		t1.Second() == t2.Second()
}

func clearCollection(t *testing.T) {
	_, err := ur.col.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		t.Fatal(err)
	}
}
func insertMockUser(t *testing.T) {
	_, err := ur.col.InsertOne(context.TODO(), u)
	if err != nil {
		t.Fatal(err)
	}
}
