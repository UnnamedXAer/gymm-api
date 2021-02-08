package usecases

import (
	"os"
	"testing"

	"github.com/unnamedxaer/gymm-api/mocks"
	"github.com/unnamedxaer/gymm-api/repositories"
)

func TestMain(m *testing.M) {
	repositories.EnsureTestEnv()

	var ur UserRepo = mocks.MockUserRepo{}

	uc = NewUserUseCases(ur)

	code := m.Run()
	os.Exit(code)
}
