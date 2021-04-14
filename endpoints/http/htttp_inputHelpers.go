package http

import (
	"github.com/unnamedxaer/gymm-api/helpers"
	"github.com/unnamedxaer/gymm-api/usecases"
)

func trimWhitespacesOnExerciseInput(e *usecases.ExerciseInput) {
	e.Name = helpers.TrimWhiteSpaces(e.Name)
	e.Description = helpers.TrimWhiteSpaces(e.Description)
}

func trimWhitespacesOnUserInput(u *usecases.UserInput) {
	u.Username = helpers.TrimWhiteSpaces(u.Username)
}
