package usecases_test

import (
	"os"
	"testing"

	"github.com/unnamedxaer/gymm-api/mocks"
	"github.com/unnamedxaer/gymm-api/testhelpers"
	"github.com/unnamedxaer/gymm-api/usecases"
)

func TestMain(m *testing.M) {
	testhelpers.EnsureTestEnv()

	var ar usecases.AuthRepo = &mocks.MockAuthRepo{}
	authUC = usecases.NewAuthUsecases(ar)

	var ur usecases.UserRepo = &mocks.MockUserRepo{}
	userUC = usecases.NewUserUseCases(ur)

	var er usecases.ExerciseRepo = &mocks.MockExerciseRepo{}
	exerciseUC = usecases.NewExerciseUseCases(er)

	var tr usecases.TrainingRepo = &mocks.MockTrainingRepo{}
	trainingUC = usecases.NewTrainingUseCases(tr)

	code := m.Run()
	os.Exit(code)
}
