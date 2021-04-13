package http

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/unnamedxaer/gymm-api/mocks"
	"github.com/unnamedxaer/gymm-api/testhelpers"
	"github.com/unnamedxaer/gymm-api/usecases"
	"github.com/unnamedxaer/gymm-api/validation"
)

var (
	app       *App
	validate  *validator.Validate
	wrongUser usecases.UserInput = usecases.UserInput{
		Username:     "1",
		EmailAddress: "email.at.no.address",
		Password:     "PWD",
	}
	correctUser usecases.UserInput = usecases.UserInput{
		Username:     "Al",
		EmailAddress: "al@mymeil.go",
		Password:     "Pwd123",
	}
)

func TestMain(m *testing.M) {
	testhelpers.EnsureTestEnv()

	validate = validation.New()
	l := &zerolog.Logger{}
	l.Level(zerolog.Disabled)
	uMockRepo := &mocks.MockUserRepo{}
	eMockRepo := &mocks.MockExerciseRepo{}
	tMockRepo := &mocks.MockTrainingRepo{}
	app = NewServer(l, uMockRepo, eMockRepo, tMockRepo, validate)
	app.AddHandlers()
	code := m.Run()
	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expectedcode, actualCode int) {
	if expectedcode != actualCode {
		t.Errorf("Expected response code %d. Got %d", expectedcode, actualCode)
	}
}
