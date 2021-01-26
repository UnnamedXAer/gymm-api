package users

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/unnamedxaer/gymm-api/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ur *UserRepository
var db *mongo.Database

func TestMain(m *testing.M) {
	err := repositories.ValidateTestEnv()
	if err != nil {
		panic(err)
	}
	db, err = repositories.GetDatabase(os.Getenv("MONGO_URI"))
	if err != nil {
		panic(err)
	}
	defer repositories.DisconnectDB(db)
	usersCol := db.Collection("users")
	ur = NewRepository(&zerolog.Logger{}, usersCol)
	code := m.Run()
	os.Exit(code)
}

func TestCreateUser(t *testing.T) {
	password, _ := hashPassword("TheSecretestPasswordEver")
	u := userData{
		FirstName:    "John",
		LastName:     "Silver",
		EmailAddress: "johnsilver@email.com",
		Password:     password,
		CreatedAt:    time.Now().UTC(),
	}
	t.Error(u)
	// ur.CreateUser()
}

func TestGetUserByID(t *testing.T) {
	password, _ := hashPassword("TheSecretestPasswordEver")
	u := userData{
		FirstName:    "John",
		LastName:     "Silver",
		EmailAddress: "johnsilver@email.com",
		Password:     password,
		CreatedAt:    time.Now().UTC(),
	}
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
		u.FirstName != gotUser.FirstName ||
		u.LastName != gotUser.LastName ||
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
