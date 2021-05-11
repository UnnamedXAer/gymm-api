package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/unnamedxaer/gymm-api/entities"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepo interface {
	GetUserByEmailAddress(ctx context.Context, emailAddress string) (*entities.AuthUser, error)
	GetUserByID(ctx context.Context, id string) (*entities.AuthUser, error)
	ChangePassword(ctx context.Context, userID string, newPwd []byte) error
	AddResetPasswordRequest(ctx context.Context, emailaddress string, expiresAt time.Time) (*entities.ResetPwdReq, error)
	SaveJWT(ctx context.Context, userID string, device string, token string, expiresAt time.Time) (*entities.UserToken, error)
	GetUserJWTs(ctx context.Context, userID string, expired entities.ExpireType) ([]entities.UserToken, error)
	DeleteJWT(ctx context.Context, ut *entities.UserToken) (int64, error)
	SaveRefreshToken(ctx context.Context, userID string, token string, expiresAt time.Time) (*entities.RefreshToken, error)
	GetRefreshToken(ctx context.Context, userID string) (*entities.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, userID string) (n int64, err error)
	DeleteRefreshTokenAndAllTokens(ctx context.Context, userID string) (n int64, err error)
}

type AuthUsecases struct {
	repo AuthRepo
	l    *zerolog.Logger
}

type IAuthUsecases interface {
	// Login checks given credentials against registered users
	Login(ctx context.Context, u *UserInput) (*entities.User, error)
	// ChangePassword updates user password
	ChangePassword(ctx context.Context, userID string, oldPwd, newPwd string) error
	// AddResetPasswordRequest adds a password reset request and send it via email
	AddResetPasswordRequest(ctx context.Context, mailer Mailer, emailaddress string) (*entities.ResetPwdReq, error)
	// SaveJWT saves jwt for given user and device name
	SaveJWT(ctx context.Context, userID string, device string, token string, expiresAt time.Time) (*entities.UserToken, error)
	// GetUserJWTs returns user jwt tokens
	// if expired is true it returns only expired tokens
	GetUserJWTs(ctx context.Context, userID string, expired entities.ExpireType) ([]entities.UserToken, error)
	// DeleteJWT removes jwt token, it returns number of deleted results and error if any
	DeleteJWT(ctx context.Context, ut *entities.UserToken) (int64, error)
	// SaveRefreshToken creates new or override existing entry in storage with given refresh token for the user.
	SaveRefreshToken(ctx context.Context, userID string, token string, expiresAt time.Time) (*entities.RefreshToken, error)
	// GetRefreshToken reads refresh token for given user.
	GetRefreshToken(ctx context.Context, userID string) (*entities.RefreshToken, error)
	// DeleteRefreshToken removes refresh token for given user id
	DeleteRefreshToken(ctx context.Context, userID string) (n int64, err error)
	// DeleteRefreshTokenAndAllTokens removes all jwt tokens and refresh token for given user
	DeleteRefreshTokenAndAllTokens(ctx context.Context, userID string) (n int64, err error)
}

type IncorrectCredentialsError struct{}

func (err IncorrectCredentialsError) Error() string {
	return "incorrect credentials"
}

func (au *AuthUsecases) Login(
	ctx context.Context,
	u *UserInput) (*entities.User, error) {
	user, err := au.repo.GetUserByEmailAddress(ctx, u.EmailAddress)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(u.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, IncorrectCredentialsError{}
		}

		return nil, errors.WithMessage(err, "usecases.Login")
	}

	return &user.User, nil
}

func (au *AuthUsecases) ChangePassword(
	ctx context.Context,
	userID string,
	oldPwd,
	newPwd string) error {
	user, err := au.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if user == nil {
		// should not happen ever
		return errors.WithMessage(
			errors.New("could not find the user"), "usecases.ChangePassword")
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(oldPwd))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return IncorrectCredentialsError{}
		}

		return errors.WithMessage(err, "usecases.ChangePassword")
	}

	passwordHash, err := hashPassword(newPwd)
	if err != nil {
		return errors.WithMessage(err, "invalid password, cannot hash")
	}

	err = au.repo.ChangePassword(ctx, userID, passwordHash)
	if err != nil {
		return err
	}

	return nil
}

func (au *AuthUsecases) AddResetPasswordRequest(
	ctx context.Context,
	mailer Mailer,
	emailaddress string) (*entities.ResetPwdReq, error) {
	expiresAt := time.Now().Add(time.Minute * 15) // @todo: config

	pwdResetReq, err := au.repo.AddResetPasswordRequest(ctx, emailaddress, expiresAt)
	if err != nil {
		var rneErr *RecordNotExistsError
		if errors.As(err, &rneErr) {
			return pwdResetReq, nil
		}
		return nil, err
	}

	go au.sendResetPwdRequestEmail(ctx, mailer, pwdResetReq)

	return pwdResetReq, nil
}

func (au *AuthUsecases) SaveJWT(
	ctx context.Context,
	userID string,
	device string,
	token string,
	expiresAt time.Time) (*entities.UserToken, error) {
	// @todo: drop previous token(s) for this device
	return au.repo.SaveJWT(ctx, userID, device, token, expiresAt)
}

func (au *AuthUsecases) DeleteJWT(
	ctx context.Context,
	token *entities.UserToken) (int64, error) {
	return au.repo.DeleteJWT(ctx, token)
}

func (au *AuthUsecases) GetUserJWTs(
	ctx context.Context,
	userID string,
	expired entities.ExpireType) ([]entities.UserToken, error) {
	return au.repo.GetUserJWTs(ctx, userID, expired)
}

func (au *AuthUsecases) SaveRefreshToken(
	ctx context.Context,
	userID string,
	token string,
	expiresAt time.Time) (*entities.RefreshToken, error) {
	return au.repo.SaveRefreshToken(ctx, userID, token, expiresAt)
}

func (au *AuthUsecases) GetRefreshToken(
	ctx context.Context,
	userID string) (*entities.RefreshToken, error) {
	return au.repo.GetRefreshToken(ctx, userID)
}

func (au *AuthUsecases) DeleteRefreshToken(
	ctx context.Context,
	userID string) (n int64, err error) {
	return au.repo.DeleteRefreshToken(ctx, userID)
}

func (au *AuthUsecases) DeleteRefreshTokenAndAllTokens(
	ctx context.Context,
	userID string) (n int64, err error) {
	return au.repo.DeleteRefreshTokenAndAllTokens(ctx, userID)
}

func (au *AuthUsecases) sendResetPwdRequestEmail(
	ctx context.Context,
	m Mailer,
	pwdResetReq *entities.ResetPwdReq) {

	if pwdResetReq == nil {
		au.l.Err(
			fmt.Errorf(
				"sending email with reset password request abandoned: nil reset password request")).Send()
		return
	}

	var data []byte
	select {
	case <-ctx.Done():
		// @todo: handle err
		return
	default:

		authUser, err := au.repo.GetUserByEmailAddress(ctx, pwdResetReq.EmailAddress)
		if err != nil {
			au.l.Err(err).Send()
			return
		}

		if authUser == nil { //should not happend
			au.l.Err(
				fmt.Errorf(
					"sending email with reset password request abandoned: user with email %s not found",
					pwdResetReq.EmailAddress)).Send()
			return
		}

		data, err = generatePwdResetEmailContent(&authUser.User, pwdResetReq)
		if err != nil {
			au.l.Err(err).Send()
			return
		}
	}

	go m.Send([]string{pwdResetReq.EmailAddress}, []byte("Reset Password"), data)
}

// NewAuthUsecases creates auth usecases
func NewAuthUsecases(l *zerolog.Logger, repo AuthRepo) IAuthUsecases {
	return &AuthUsecases{
		repo: repo,
		l:    l,
	}
}
