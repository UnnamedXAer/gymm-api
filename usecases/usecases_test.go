package usecases

import (
	"os"
	"testing"

	"github.com/unnamedxaer/gymm-api/mocks"
	"github.com/unnamedxaer/gymm-api/testhelpers"
)

func TestMain(m *testing.M) {
	testhelpers.EnsureTestEnv()

	var ur UserRepo = mocks.MockUserRepo{}

	uc = NewUserUseCases(ur)

	code := m.Run()
	os.Exit(code)
}
