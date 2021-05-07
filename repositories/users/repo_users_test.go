package users

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/unnamedxaer/gymm-api/repositories"
	"github.com/unnamedxaer/gymm-api/testhelpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ur *UserRepository
	// db *mongo.Database
	u UserData
)

func TestMain(m *testing.M) {
	logger := zerolog.New(os.Stdout)
	testhelpers.EnsureTestEnv()

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
	defer testhelpers.DisconnectDB(&logger, db)

	usersCol := db.Collection("users")
	_, err = usersCol.DeleteMany(context.Background(), bson.D{})
	if err != nil {
		panic(err)
	}

	ur = NewRepository(&zerolog.Logger{}, usersCol)

	// password := []byte("TheSecretestPasswordEver123$%^")
	u = UserData{
		Username:     "John Silver",
		EmailAddress: "johnsilver@email.com",
		Password:     []byte("TheSecretestPasswordEver123$%^"),
		CreatedAt:    time.Now().UTC(),
	}

	code := m.Run()
	os.Exit(code)
}

func TestCreateUser(t *testing.T) {
	ctx := context.TODO()
	clearCollection(t)
	gotUser, err := ur.CreateUser(ctx, u.Username, u.EmailAddress, u.Password)
	if err != nil {
		t.Fatal(err)
	}

	if gotUser == nil {
		t.Fatalf("want to get user base on data: %v, got nil", u)
	}

	if u.EmailAddress != gotUser.EmailAddress ||
		u.Username != gotUser.Username {
		t.Fatalf("want user base on data: %v, got: %v",
			u,
			gotUser)
	}

	if gotUser.ID == "" {
		t.Fatalf("want 'ID' to NOT be zero value, got %q for email: %q",
			gotUser.ID,
			u.EmailAddress)
	}

	if gotUser.CreatedAt.IsZero() {
		t.Fatalf("want 'CreateAt' to NOT be zero value, got %q for email: %q",
			gotUser.CreatedAt,
			u.EmailAddress)
	}
}

func TestCreateUserDuplicatedEmail(t *testing.T) {
	ctx := context.TODO()
	clearCollection(t)
	insertMockUser(t)
	gotUser, err := ur.CreateUser(ctx, u.Username, u.EmailAddress, u.Password)
	if err != nil {
		if errors.Is(err, repositories.NewErrorEmailAddressInUse()) {
			return
		}
		t.Fatal(err)
	}

	t.Fatalf("want error %q for email addres: %q, got %v",
		repositories.NewErrorEmailAddressInUse(),
		u.EmailAddress,
		gotUser)
}

func TestGetUserByID(t *testing.T) {
	ctx := context.TODO()
	clearCollection(t)
	results, err := ur.col.InsertOne(context.TODO(), u)
	if err != nil {
		t.Fatal(err)
	}

	uID := results.InsertedID.(primitive.ObjectID).Hex()

	gotUser, err := ur.GetUserByID(ctx, uID)
	if err != nil {
		t.Fatalf("want user, got %v", err)
	}

	if gotUser == nil || (testhelpers.TimesEqual(u.CreatedAt, gotUser.CreatedAt) == false ||
		u.EmailAddress != gotUser.EmailAddress ||
		u.Username != gotUser.Username ||
		uID != gotUser.ID) {
		t.Errorf("want user like: %v, got: %v", u, gotUser)
	}
}

func TestGetUserByIDNotExisting(t *testing.T) {
	ctx := context.TODO()
	clearCollection(t)
	uID := "60108393da81e60598d5347f"

	gotUser, err := ur.GetUserByID(ctx, uID)
	if err != nil {
		t.Fatalf("want nil error, got: %v", err)
	}

	if gotUser != nil {
		t.Fatalf("want nil user, got %v", gotUser)
	}
}

func TestGetUserByIDInvalidID(t *testing.T) {
	ctx := context.TODO()
	clearCollection(t)
	uID := "6s108393da81e60598d5347f"

	gotUser, err := ur.GetUserByID(ctx, uID)
	var e *repositories.InvalidIDError
	if !errors.As(err, &e) {
		t.Fatalf("want error: %q, got: %v", repositories.NewErrorInvalidID(uID, "user"), err)
	}

	if gotUser != nil {
		t.Fatalf("want nil user, got %v", err)
	}
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
