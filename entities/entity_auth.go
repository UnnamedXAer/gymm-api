package entities

type AuthUser struct {
	User
	Password []byte
}
