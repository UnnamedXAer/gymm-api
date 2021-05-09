package mocks

var (
	UserID            = "6072d3206144644984a54fa1"
	NonexistingUserID = UserID[:len(UserID)-1] + "a"

	NonexistingEmail = "notfound@example.com"
	Password         = []byte("TheSecretestPasswordEver123$%^")
)
