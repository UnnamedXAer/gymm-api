package exercises

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
	"github.com/unnamedxaer/gymm-api/mocks"
	"github.com/unnamedxaer/gymm-api/repositories"
	"github.com/unnamedxaer/gymm-api/testhelpers"
	"github.com/unnamedxaer/gymm-api/usecases"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	exerciseColName = "exercises"
)

var (
	exerciseRepo   usecases.ExerciseRepo
	mockedExercise entities.Exercise
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

	exercisesCol := db.Collection(exerciseColName)
	res, err := exercisesCol.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		log.Fatalln(err)
	}

	exerciseRepo = NewRepository(&loggerMock, exercisesCol)

	if err != nil {
		loggerMock.Err(err).Msgf("%d", res.DeletedCount)
		panic(err)
	}
	me, err := mocks.InsertMockExercise(exerciseRepo)
	if err != nil {
		panic(err)
	}
	mockedExercise = *me

	os.Exit(m.Run())
}

func TestGetExerciseByID(t *testing.T) {
	ctx := context.TODO()
	want := mockedExercise
	ex, err := exerciseRepo.GetExerciseByID(ctx, want.ID)
	if err != nil {
		t.Error(err)
		return
	}

	if ex.ID != want.ID {
		t.Errorf("want 'ID' to be %q, got %q", want.ID, ex.ID)
	}

	if ex.Name != want.Name {
		t.Errorf("want 'Name' to be %q, got %q", want.Name, ex.Name)
	}

	if ex.Description != want.Description {
		t.Errorf("want 'Description' to be %q, got %q", want.Description, ex.Description)
	}

	if !testhelpers.TimesEqual(ex.CreatedAt, want.CreatedAt) {
		t.Errorf("want 'CreateAt' to be %s, got %s", want.CreatedAt, ex.CreatedAt)
	}

	if ex.CreatedBy != want.CreatedBy {
		t.Errorf("want 'CreateBy' to be %q, got %q", want.CreatedBy, ex.CreatedBy)
	}

	if ex.SetUnit != want.SetUnit {
		t.Errorf("want 'SetUnit' to be Time (%d), got %d", want.SetUnit, ex.SetUnit)
	}
	mockedExercise = *ex
}

func TestCreateExercise(t *testing.T) {
	ctx := context.TODO()
	want := mockedExercise
	want.Name += fmt.Sprintf("-> %d", time.Now().UnixNano())
	ex, err := exerciseRepo.CreateExercise(ctx, want.Name, want.Description, want.SetUnit, want.CreatedBy)
	if err != nil {
		t.Error(err)
		return
	}

	if ex.ID == "" {
		t.Error("want 'ID' not to be zero value")
		return
	}

	if ex.Name != want.Name {
		t.Errorf("want 'Name' to be %q, got %q", want.Name, ex.Name)
	}

	if ex.Description != want.Description {
		t.Errorf("want 'Description' to be %q, got %q", want.Description, ex.Description)
	}

	if ex.CreatedAt.IsZero() {
		t.Error("want 'CreateAt' NOT to be zero value")
	}

	if ex.CreatedBy != want.CreatedBy {
		t.Errorf("want 'CreateBy' to be %q, got %q", want.CreatedBy, ex.CreatedBy)
	}

	if ex.SetUnit != want.SetUnit {
		t.Errorf("want 'SetUnit' to be Time (%d), got %d", want.SetUnit, ex.SetUnit)
	}
	mockedExercise = *ex
}

func TestUpdateExercise(t *testing.T) {
	ctx := context.TODO()
	want := mockedExercise
	want.Description += "\n-> updated at " + time.Now().String()
	want.SetUnit = entities.Time
	ex, err := exerciseRepo.UpdateExercise(ctx, &want)
	if err != nil {
		t.Error(err)
		return
	}

	if ex.ID != want.ID {
		t.Errorf("want 'ID' to be %q, got %q", want.ID, ex.ID)
	}

	if ex.Name != want.Name {
		t.Errorf("want 'Name' to be %q, got %q", want.Name, ex.Name)
	}

	if ex.Description != want.Description {
		t.Errorf("want 'Description' to be %q, got %q", want.Description, ex.Description)
	}

	if !testhelpers.TimesEqual(ex.CreatedAt, want.CreatedAt) {
		t.Errorf("want 'CreateAt' to be %s, got %s", want.CreatedAt, ex.CreatedAt)
	}

	if ex.CreatedBy != want.CreatedBy {
		t.Errorf("want 'CreateBy' to be %q, got %q", want.CreatedBy, ex.CreatedBy)
	}

	if ex.SetUnit != want.SetUnit {
		t.Errorf("want 'SetUnit' to be Time (%d), got %d", want.SetUnit, ex.SetUnit)
	}
	mockedExercise = *ex
}

func TestUpdateExerciseOneProp(t *testing.T) {
	ctx := context.TODO()
	input := entities.Exercise{
		ID:      mockedExercise.ID,
		SetUnit: entities.Time,
	}
	want := mockedExercise
	want.SetUnit = input.SetUnit
	ex, err := exerciseRepo.UpdateExercise(ctx, &input)
	if err != nil {
		t.Error(err)
		return
	}

	if ex.ID != want.ID {
		t.Errorf("want 'ID' to be %q, got %q", want.ID, ex.ID)
	}

	if ex.Name != want.Name {
		t.Errorf("want 'Name' to be %q, got %q", want.Name, ex.Name)
	}

	if ex.Description != want.Description {
		t.Errorf("want 'Description' to be %q, got %q", want.Description, ex.Description)
	}

	if !testhelpers.TimesEqual(ex.CreatedAt, want.CreatedAt) {
		t.Errorf("want 'CreateAt' to be %s, got %s", want.CreatedAt, ex.CreatedAt)
	}

	if ex.CreatedBy != want.CreatedBy {
		t.Errorf("want 'CreateBy' to be %q, got %q", want.CreatedBy, ex.CreatedBy)
	}

	if ex.SetUnit != want.SetUnit {
		t.Errorf("want 'SetUnit' to be Time (%d), got %d", want.SetUnit, ex.SetUnit)
	}
	mockedExercise = *ex
}

func TestGetExercisesByName(t *testing.T) {
	ctx := context.TODO()
	want := mockedExercise
	name := strings.ToLower(want.Name[:len(mockedExercise.Name)-1])
	exercises, err := exerciseRepo.GetExercisesByName(ctx, name)
	if err != nil {
		t.Error(err)
		return
	}

	var found bool
	for _, ex := range exercises {
		if ex.ID == want.ID {
			found = true
		}
	}

	if !found {
		t.Errorf("want find exercise with ID %q for name %q, got %v", want.ID, name, exercises)
	}
}

func TestGetExercisesByNameNotExisting(t *testing.T) {
	ctx := context.TODO()
	name := "notfound"
	exercises, err := exerciseRepo.GetExercisesByName(ctx, name)
	if err != nil {
		t.Error(err)
		return
	}

	if len(exercises) > 0 {
		t.Errorf("want get 0 exercises, got %v", exercises)
	}
}
