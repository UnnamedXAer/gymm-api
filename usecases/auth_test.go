package usecases_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/unnamedxaer/gymm-api/mailer"
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

func TestAddResetPasswordRequest(t *testing.T) {
	testCases := []struct {
		desc         string
		emailAddress string
		errTxt       string
	}{
		{
			desc:   "missing email",
			errTxt: "",
		},
		{
			desc:         "nonexisting user",
			emailAddress: mocks.NonexistingEmail,
			errTxt:       "",
		},
		{
			desc:         "correct",
			emailAddress: mocks.ExampleUser.EmailAddress,
			errTxt:       "",
		},
		{
			desc:         "correct2",
			emailAddress: "unnamedxaer@gmail.com",
			errTxt:       "",
		},
	}

	// mockedMailer := &mocks.MockMailer{}
	mockedMailer := mailer.NewMailer(
		&mockedLogger,
		func(err error) {
			fmt.Printf("mailer-mock: %v", err)
		})

	ctx := context.TODO()
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			pwdResReq, err := authUC.AddResetPasswordRequest(ctx, mockedMailer, tC.emailAddress)
			if tC.errTxt == "" {
				if err != nil {
					t.Errorf("want nil error, got %q", err)
				}
			} else {
				if err == nil || !strings.Contains(err.Error(), tC.errTxt) {
					t.Errorf("want error like %q, got %q", tC.errTxt, err)
				}
			}

			if err != nil && pwdResReq == nil && !strings.Contains(err.Error(), usecases.NewErrorRecordNotExists("user").Error()) {
				t.Error("want saved request, got nil")
			}

		})
	}

}
