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

	jwtKey := []byte(os.Getenv("JWT_KEY"))
	if len(jwtKey) < 10 {
		l.Panic().Msg("missing or too short jwt key")
	}

	aMockRepo := &mocks.MockAuthRepo{}
	uMockRepo := &mocks.MockUserRepo{}
	eMockRepo := &mocks.MockExerciseRepo{}
	tMockRepo := &mocks.MockTrainingRepo{}
	app = NewServer(l, aMockRepo, uMockRepo, eMockRepo, tMockRepo, validate, jwtKey)
	app.AddHandlers()
	code := m.Run()
	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, want, got int) {
	if want != got {
		t.Errorf("want response code %d, got %d", want, got)
	}
}
