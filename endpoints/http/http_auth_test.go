package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/unnamedxaer/gymm-api/entities"
	"github.com/unnamedxaer/gymm-api/mocks"
	"github.com/unnamedxaer/gymm-api/usecases"
)

func TestLogin(t *testing.T) {
	u := usecases.UserInput{
		EmailAddress: mocks.ExampleUser.EmailAddress,
		Password:     string(mocks.Password),
	}

	uJSON, err := json.Marshal(&u)
	if err != nil {
		t.Fatalf("could not marshal login payload: %v", err)
	}

	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(uJSON))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	res := response.Result()
	cookies := res.Cookies()

	var tokenCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == cookieJwtTokenName {
			tokenCookie = c
		}
	}

	if tokenCookie == nil {
		t.Errorf("want token cookie, got nil")
	} else {

		if !tokenCookie.HttpOnly {
			t.Errorf("want token cookie to be HttpOnly")
		}

		if !tokenCookie.Expires.IsZero() {
			t.Errorf("want token cookie expire time to be zero, got %v", tokenCookie.Expires)
		}
	}

	var got struct {
		Error string
		User  *entities.User
	}
	err = json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Fatalf("could not decode login response: %v", err)
	}

	if got.User == nil || got.Error != "" {
		t.Errorf("want authenticated user and no error, got %#v", got)
	}

}

func TestRegister(t *testing.T) {
	u := correctUser

	uJSON, _ := json.Marshal(u)

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(uJSON))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)
}

func TestRegisterMalformedData(t *testing.T) {
	uStr := `{"userName: "Al", "emailAddress" : "al@mymeil.com", "password":"pwd123"}`
	uJSON := []byte(uStr)

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(uJSON))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusUnprocessableEntity, response.Code)
}

func TestRegisterValidationFail(t *testing.T) {
	u := wrongUser

	uJSON, _ := json.Marshal(u)

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(uJSON))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotAcceptable, response.Code)
}

func TestChangePassword(t *testing.T) {
	newPassword := string(mocks.Password) + "X"

	testCases := []struct {
		desc   string
		oldPwd string
		newPwd string
		code   int
		errTxt string
	}{
		{
			desc:   "missing old password",
			newPwd: newPassword,
			code:   http.StatusUnauthorized,
			errTxt: "incorrect",
		},
		{
			desc:   "incorrect old password",
			oldPwd: correctUser.Password + "👽",
			newPwd: newPassword,
			code:   http.StatusUnauthorized,
			errTxt: "incorrect",
		},
		{
			desc:   "missing new password",
			oldPwd: correctUser.Password,
			code:   http.StatusBadRequest,
			errTxt: "'password' field value is required",
		},
		{
			desc:   "correct",
			oldPwd: string(mocks.Password),
			newPwd: newPassword,
			code:   http.StatusOK,
			errTxt: "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			payload := map[string]string{
				"oldPassword": tC.oldPwd,
				"password":    tC.newPwd,
			}

			b, _ := json.Marshal(&payload)

			req, _ := http.NewRequest(http.MethodPost, "/password/change", bytes.NewBuffer(b))
			req.Header.Set("Content-Type", "application/json")

			res := executeRequest(req)

			checkResponseCode(t, tC.code, res.Code)

			got := res.Body.String()

			if tC.errTxt == "" {
				if got != "" {
					t.Errorf("want empty body, got %q", got)
				}
			}
			if tC.errTxt != "" && !strings.Contains(got, tC.errTxt) {
				t.Errorf("want error like %q, got %q", tC.errTxt, got)
			}

		})
	}
}

func TestAddResetPasswordRequest(t *testing.T) {
	testCases := []struct {
		desc         string
		emailAddress string
		errTxt       string
		code         int
	}{
		{
			desc:   "missing email",
			errTxt: "'emailAddress' field value is required",
			code:   http.StatusBadRequest,
		},
		{
			desc:         "nonexisting user",
			emailAddress: mocks.NonexistingEmail,
			errTxt:       "",
			code:         http.StatusAccepted,
		},
		{
			desc:         "correct",
			emailAddress: mocks.ExampleUser.EmailAddress,
			errTxt:       "",
			code:         http.StatusAccepted,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			payload := map[string]string{
				"emailAddress": tC.emailAddress,
			}

			b, _ := json.Marshal(&payload)

			req, _ := http.NewRequest(http.MethodPost, "/password/reset", bytes.NewBuffer(b))
			req.Header.Set("Content-Type", "application/json")

			res := executeRequest(req)

			checkResponseCode(t, tC.code, res.Code)

			got := res.Body.String()

			if tC.errTxt == "" && len(got) != 0 {
				t.Errorf("want empty body, got %q", got)
				return
			}

			if !strings.Contains(got, tC.errTxt) {
				t.Errorf("want error like %q, got %q", tC.errTxt, got)
			}

		})
	}
}
