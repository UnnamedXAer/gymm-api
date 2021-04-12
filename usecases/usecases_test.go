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

	var ur usecases.UserRepo = mocks.MockUserRepo{}
	userUC = usecases.NewUserUseCases(ur)

	var er usecases.ExerciseRepo = &mocks.MockExerciseRepo{}
	exerciseUC = usecases.NewExerciseUseCases(er)

	code := m.Run()
	os.Exit(code)
}
