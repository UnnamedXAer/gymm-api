package auth

import (
	"context"
	"log"
	"os"
	"strings"
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
	nonexistingID = mocks.UserID[:len(mocks.UserID)-1] + "a"
	pwdHash       []byte

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

	if mocks.UserID[len(mocks.UserID)-1] == 'a' {
		nonexistingID = nonexistingID[:len(nonexistingID)-1] + "b"
	}

	result := usersCol.FindOne(context.TODO(),
		bson.M{
			"$or": bson.A{
				bson.M{"email_address": mocks.ExampleUser.EmailAddress},
				bson.M{"_id": nonexistingEmail},
			},
		})
	if err = result.Err(); err != nil {
		if err != mongo.ErrNoDocuments {
			log.Fatalln(err)
		}
		pwdHash, err = bcrypt.GenerateFromPassword([]byte(mocks.Password), bcrypt.MinCost)
		if err != nil {
			log.Fatalln(err)
		}
		usersRepo := users.NewRepository(&logger, usersCol)
		u, err := usersRepo.CreateUser(
			context.TODO(),
			mocks.ExampleUser.Username,
			mocks.ExampleUser.EmailAddress,
			pwdHash,
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
	testCases := []struct {
		desc        string
		id          string
		errTxt      string
		returnsUser bool
	}{
		{
			desc:        "invalid id",
			id:          mocks.UserID + "🐼",
			errTxt:      repositories.NewErrorInvalidID(mocks.UserID+"🐼", "user").Error(),
			returnsUser: false,
		},
		{
			desc:        "not existing user",
			id:          mocks.UserID,
			errTxt:      "",
			returnsUser: false,
		},
		{
			desc:        "correct id",
			id:          mockedUser.ID,
			errTxt:      "",
			returnsUser: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctx := context.TODO()
			got, err := authRepo.GetUserByID(ctx, tC.id)

			if tC.returnsUser && got == nil {
				t.Errorf("want user, got nil")
			}
			if !tC.returnsUser && got != nil {
				t.Errorf("want nil user , got %v", got)
			}

			if tC.errTxt == "" {
				if err != nil {
					t.Errorf("want nil error, got %v", err)
				}
			} else {
				if !strings.Contains(err.Error(), tC.errTxt) {
					t.Errorf("want error %q, got %v", tC.errTxt, err)
				}
			}
		})
	}
}

func TestChangePassword(t *testing.T) {
	testCases := []struct {
		desc   string
		id     string
		pwd    []byte
		errTxt string
	}{
		{
			desc:   "invalid id",
			id:     mocks.UserID + "🐼",
			pwd:    pwdHash,
			errTxt: repositories.NewErrorInvalidID(mocks.UserID+"🐼", "user").Error(),
		},
		{
			desc:   "not existing user",
			id:     nonexistingID,
			pwd:    pwdHash,
			errTxt: "no record has been updated",
		},
		{
			desc:   "new password the same as old",
			id:     mockedUser.ID,
			pwd:    pwdHash,
			errTxt: "",
		},
		{
			desc:   "new password different than old",
			id:     mockedUser.ID,
			pwd:    append(pwdHash, 'a', 'b'),
			errTxt: "",
		},
		{
			desc:   "__restore default password",
			id:     mockedUser.ID,
			pwd:    pwdHash,
			errTxt: "",
		},
	}
	var err error

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctx := context.TODO()
			err = authRepo.ChangePassword(ctx, tC.id, tC.pwd)

			if tC.errTxt == "" {
				if err != nil {
					t.Errorf("want nil error, got %v", err)
				}
			} else {
				if !strings.Contains(err.Error(), tC.errTxt) {
					t.Errorf("want error %q, got %v", tC.errTxt, err)
				}
			}
		})
	}
}

func TestGetUserByEmailAddress(t *testing.T) {
	ctx := context.TODO()
	got, err := authRepo.GetUserByEmailAddress(ctx, mockedUser.EmailAddress)
	if err != nil {
		t.Fatalf("want user, got error: %v", err)
	}

	if got == nil || (testhelpers.TimesEqual(mockedUser.CreatedAt, got.CreatedAt) == false ||
		mockedUser.EmailAddress != got.EmailAddress ||
		mockedUser.Username != got.Username ||
		mockedUser.ID != got.ID ||
		len(mockedUser.Password) != len(got.Password)) {
		t.Errorf("want user like: %v, got: %v", mockedUser, got)
	}
}

func TestGetUserByEmailAddressNotExists(t *testing.T) {
	ctx := context.TODO()
	got, err := authRepo.GetUserByEmailAddress(ctx, nonexistingEmail)
	if err != nil {
		t.Fatalf("want nil error, got: %v", err)
	}

	if got != nil {
		t.Errorf("want nil, got: %v", got)
	}
}

func TestGetUserByEmailAddressEmpty(t *testing.T) {
	ctx := context.TODO()
	got, err := authRepo.GetUserByEmailAddress(ctx, "")
	if err == nil {
		t.Fatal("want error 'empty email address', got: nil")
	}
	if !strings.Contains(err.Error(), "empty email address") {
		t.Fatalf("want error 'empty email address', got: %v", err)
	}

	if got != nil {
		t.Errorf("want nil, got: %v", got)
	}
}
