package usecases_test

import (
	"context"
	"strings"
	"testing"

	"github.com/unnamedxaer/gymm-api/mocks"
	"github.com/unnamedxaer/gymm-api/usecases"
)

var (
	authUC usecases.IAuthUsecases
)

func TestChangePassword(t *testing.T) {
	newPassword := string(mocks.Password) + "X"
	testCases := []struct {
		desc   string
		oldPwd string
		newPwd string
		errTxt string
	}{
		{
			desc:   "missing old password",
			newPwd: newPassword,
			errTxt: "incorrect credentials",
		},
		{
			desc:   "incorrect old password",
			oldPwd: string(mocks.Password) + "👽",
			newPwd: newPassword,
			errTxt: "incorrect credentials",
		},
		{
			desc:   "correct",
			oldPwd: string(mocks.Password),
			newPwd: newPassword,
			errTxt: "",
		},
	}

	ctx := context.TODO()
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			err := authUC.ChangePassword(ctx, mocks.UserID, tC.oldPwd, tC.newPwd)
			if tC.errTxt == "" {
				if err != nil {
					t.Errorf("want nil error, got %q", err)
				}
			} else {
				if err == nil || !strings.Contains(err.Error(), tC.errTxt) {
					t.Errorf("want error like %q, got %q", tC.errTxt, err)
				}
			}
		})
	}
}

func TestResetPassword(t *testing.T) {
	testCases := []struct {
		desc         string
		emailAddress string
		errTxt       string
		code         int
	}{

		{
			desc:   "missing email",
			errTxt: "'emailAddress' field value is required",
		},
		{
			desc:         "nonexisting user",
			emailAddress: mocks.NonexistingEmail,
			errTxt:       "no documents",
		},
		{
			desc:         "correct",
			emailAddress: mocks.ExampleUser.EmailAddress,
			errTxt:       "",
		},
	}

	ctx := context.TODO()
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			pwdResReq, err := authUC.AddResetPasswordRequest(ctx, tC.emailAddress)
			if tC.errTxt == "" {
				if err != nil {
					t.Errorf("want nil error, got %q", err)
				}
			} else {
				if err == nil || !strings.Contains(err.Error(), tC.errTxt) {
					t.Errorf("want error like %q, got %q", tC.errTxt, err)
				}
			}

			if len(pwdResReq.ID) == 0 {
				t.Errorf("want saved request, got %v", pwdResReq)
			}

			if pwdResReq.EmailAddress != tC.emailAddress {
				t.Errorf("want request for %s, got %v", tC.emailAddress, pwdResReq)
			}
		})
	}
}
