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

	uc = usecases.NewUserUseCases(ur)

	code := m.Run()
	os.Exit(code)
}
