package auth

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/mocks"
	"github.com/unnamedxaer/gymm-api/repositories"
	"github.com/unnamedxaer/gymm-api/repositories/users"
	"github.com/unnamedxaer/gymm-api/testhelpers"
	"github.com/unnamedxaer/gymm-api/usecases"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const (
	usersCollectionName  = "users"
	tokensCollectionName = "tokens"
	nonexistingEmail     = "notfound@example.com"
)

var (
	authRepo   usecases.AuthRepo
	mockedUser entities.AuthUser
)

func TestMain(m *testing.M) {

	testhelpers.EnsureTestEnv()
	logger := zerolog.New(os.Stdout)

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatalln("environment variable 'DB_NAME' is not set")
	}
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatalln("environment variable 'MONGO_URI' is not set")
	}
	db, err := repositories.GetDatabase(&logger, mongoURI, dbName)
	if err != nil {
		log.Fatalln(err)
	}

	err = repositories.CreateCollections(&logger, db)
	if err != nil {
		log.Fatalln(err)
	}
	defer testhelpers.DisconnectDB(&logger, db)

	tokensCol := db.Collection(usersCollectionName)
	refTokensCol := db.Collection(usersCollectionName)
	usersCol := db.Collection(usersCollectionName)
	_, err = usersCol.DeleteOne(context.TODO(), bson.M{"email_address": nonexistingEmail})
	if err != nil && err != mongo.ErrNoDocuments {
		log.Fatalln(err)
	}

	result := usersCol.FindOne(context.TODO(), bson.M{"email_address": mocks.ExampleUser.EmailAddress})
	if err = result.Err(); err != nil {
		if err != mongo.ErrNoDocuments {
			log.Fatalln(err)
		}
		hashed, err := bcrypt.GenerateFromPassword([]byte(mocks.Password), bcrypt.MinCost)
		if err != nil {
			log.Fatalln(err)
		}
		usersRepo := users.NewRepository(&logger, usersCol)
		u, err := usersRepo.CreateUser(
			mocks.ExampleUser.Username,
			mocks.ExampleUser.EmailAddress,
			hashed,
		)
		if err != nil {
			log.Fatalln(err)
		}
		mockedUser = entities.AuthUser{
			User:     *u,
			Password: mocks.Password,
		}
	} else {
		data := users.UserData{}
		err := result.Decode(&data)
		if err != nil {
			log.Fatalln(err)
		}
		mockedUser.User = entities.User{
			ID:           data.ID.Hex(),
			EmailAddress: data.EmailAddress,
			Username:     data.Username,
			CreatedAt:    data.CreatedAt,
		}
		mockedUser.Password = data.Password
	}

	authRepo = NewRepository(&logger, usersCol, tokensCol, refTokensCol)

	os.Exit(m.Run())
}

func TestGetUserByID(t *testing.T) {
	got, err := authRepo.GetUserByEmailAddress(mockedUser.EmailAddress)
	if err != nil {
		t.Fatalf("want user, got error: %v", err)
	}

	if got == nil || (testhelpers.TimesEqual(mockedUser.CreatedAt, got.CreatedAt) == false ||
		mockedUser.EmailAddress != got.EmailAddress ||
		mockedUser.Username != got.Username ||
		mockedUser.ID != got.ID ||
		len(got.Password) == 0) {
		t.Errorf("want user like: %v, got: %v", mockedUser, got)
	}
}

func TestGetUserByIDNotExists(t *testing.T) {
	got, err := authRepo.GetUserByEmailAddress(nonexistingEmail)
	if err != nil {
		t.Fatalf("want nil error, got: %v", err)
	}

	if got != nil {
		t.Errorf("want nil, got: %v", got)
	}
}
