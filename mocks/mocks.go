package mocks

import "time"

var (
	UserID            = "6072d3206144644984a54fa1"
	NonexistingUserID = UserID[:len(UserID)-1] + "a"

	NonexistingEmail = "notfound@example.com"
	Password         = []byte("TheSecretestPasswordEver123$%^")
	PasswordHash     = []byte("$2a$04$d0sgKcu9y.h8grIpktLj9OAdcv7pGy5CZ9aaz5zqPAkPyqlxLGF5W")

	Now = time.Now().UTC()
)
