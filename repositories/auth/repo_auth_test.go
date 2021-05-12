package auth

import (
	"context"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/mocks"
	"github.com/unnamedxaer/gymm-api/repositories"
	"github.com/unnamedxaer/gymm-api/repositories/users"
	"github.com/unnamedxaer/gymm-api/testhelpers"
	"github.com/unnamedxaer/gymm-api/usecases"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	authRepo   usecases.AuthRepo
	mockedUser entities.AuthUser
)

func TestMain(m *testing.M) {

	testhelpers.EnsureTestEnv()
	loggerMock := zerolog.New(nil)

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatalln("environment variable 'DB_NAME' is not set")
	}
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatalln("environment variable 'MONGO_URI' is not set")
	}
	db, err := repositories.GetDatabase(&loggerMock, mongoURI, dbName)
	if err != nil {
		log.Fatalln(err)
	}

	err = repositories.CreateCollections(&loggerMock, db)
	if err != nil {
		log.Fatalln(err)
	}
	defer testhelpers.DisconnectDB(&loggerMock, db)

	tokensCol := db.Collection(repositories.TokensCollectionName)
	refTokensCol := db.Collection(repositories.RefreshTokensCollectionName)
	pwdResReqCol := db.Collection(repositories.ResPwdReqCollectionName)
	usersCol := db.Collection(repositories.UsersCollectionName)
	_, err = usersCol.DeleteOne(context.TODO(), bson.M{"email_address": mocks.NonexistingEmail})
	if err != nil && err != mongo.ErrNoDocuments {
		log.Fatalln(err)
	}

	if mocks.UserID[len(mocks.UserID)-1] == 'a' {
		mocks.NonexistingUserID = mocks.NonexistingUserID[:len(mocks.NonexistingUserID)-1] + "b"
	}

	update := bson.M{"$set": users.UserData{
		Username:     mocks.ExampleUser.Username,
		EmailAddress: mocks.ExampleUser.EmailAddress,
		Password:     mocks.PasswordHash,
		CreatedAt:    time.Now(),
	}}

	result := usersCol.FindOneAndUpdate(context.TODO(),
		bson.M{"email_address": mocks.ExampleUser.EmailAddress}, update,
		options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true))
	if err = result.Err(); err != nil {
		if err != mongo.ErrNoDocuments {
			log.Fatalln(err)
		}
		usersRepo := users.NewRepository(&loggerMock, usersCol)
		u, err := usersRepo.CreateUser(
			context.TODO(),
			mocks.ExampleUser.Username,
			mocks.ExampleUser.EmailAddress,
			mocks.PasswordHash,
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

	authRepo = NewRepository(&loggerMock, usersCol, tokensCol, refTokensCol, pwdResReqCol)

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
			errTxt:      usecases.NewErrorInvalidID(mocks.UserID+"🐼", "user").Error(),
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

	pwdHash := make([]byte, len(mocks.PasswordHash))

	copy(pwdHash, mocks.PasswordHash)

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
			errTxt: usecases.NewErrorInvalidID(mocks.UserID+"🐼", "user").Error(),
		},
		{
			desc:   "not existing user",
			id:     mocks.NonexistingUserID,
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
			pwd:    mocks.PasswordHash,
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

func TestAddResetPasswordRequest(t *testing.T) {
	testCases := []struct {
		desc         string
		emailAddress string
		expiresAt    time.Time
		errTxt       string
		code         int
	}{
		{
			desc:      "missing email",
			expiresAt: time.Now().Add(time.Minute * 15),
			errTxt:    usecases.NewErrorRecordNotExists("user").Error(),
		},
		{
			desc:         "nonexisting user",
			expiresAt:    time.Now().Add(time.Minute * 15),
			emailAddress: mocks.NonexistingEmail,
			errTxt:       usecases.NewErrorRecordNotExists("user").Error(),
		},
		{
			desc:         "past expiration time",
			expiresAt:    time.Now().Add(time.Minute * -15),
			emailAddress: mocks.NonexistingEmail,
			errTxt:       "expiration time from the past",
		},
		{
			desc:         "zero value expiration time",
			expiresAt:    time.Time{},
			emailAddress: mocks.NonexistingEmail,
			errTxt:       "expiration time from the past",
		},
		{
			desc:         "correct",
			expiresAt:    time.Now().Add(time.Minute * 15),
			emailAddress: mocks.ExampleUser.EmailAddress,
			errTxt:       "",
		},
	}

	var err error
	var pwdResReq *entities.ResetPwdReq

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctx := context.TODO()
			pwdResReq, err = authRepo.AddResetPasswordRequest(ctx, tC.emailAddress, tC.expiresAt)

			if tC.errTxt == "" {
				if err != nil {
					t.Errorf("want nil error, got %v", err)
				}
			} else {
				if pwdResReq != nil {
					t.Errorf("want nil pwd request, got %v", pwdResReq)
				}
				if !strings.Contains(err.Error(), tC.errTxt) {
					t.Errorf("want error %q, got %v", tC.errTxt, err)
				}
			}

			if pwdResReq == nil {
				if tC.errTxt == "" {
					t.Errorf("want password reset request, got nil")
				}
				return
			}

			if len(pwdResReq.ID) == 0 {
				t.Errorf("want saved request, got %v", pwdResReq)
			}

			if pwdResReq.EmailAddress != tC.emailAddress {
				t.Errorf("want request for %s, got %v", tC.emailAddress, pwdResReq)
			}
		})
	}
}

func TestUpdatePasswordForResetRequest(t *testing.T) {

	nonexistingReqID := mocks.ExampleResetPwdReq.ID[:len(mocks.ExampleResetPwdReq.ID)-1]
	if mocks.ExampleResetPwdReq.ID[len(mocks.ExampleResetPwdReq.ID)-1] == 'f' {
		nonexistingReqID += "e"
	} else {
		nonexistingReqID += "f"
	}

	testCases := []struct {
		desc     string
		password []byte
		reqID    string
		errTxt   string
	}{
		{
			desc:   "missing pasword",
			reqID:  mocks.ExampleResetPwdReq.ID,
			errTxt: "missing password",
		},
		{
			desc:     "missing id",
			password: pwdHash,
			errTxt:   usecases.NewErrorInvalidID("", "reset password request").Error(),
		},
		{
			desc:     "not existing request",
			password: pwdHash,
			reqID:    "notfound",
			errTxt:   usecases.NewErrorRecordNotExists("reset password request").Error(),
		},
		{
			desc:     "correct",
			reqID:    nonexistingReqID,
			password: pwdHash,
		},
	}

	ctx := context.TODO()
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			err := authRepo.UpdatePasswordForResetRequest(ctx, tC.reqID, tC.password)
			if tC.errTxt == "" {
				if err != nil {
					t.Errorf("want nil error, got %q", err)
				}
			} else {
				if err == nil || !strings.Contains(err.Error(), tC.errTxt) {
					t.Errorf("want error like %q, got %q", tC.errTxt, err)
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
	got, err := authRepo.GetUserByEmailAddress(ctx, mocks.NonexistingEmail)
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
