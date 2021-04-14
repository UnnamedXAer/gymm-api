package http

import (
	"github.com/unnamedxaer/gymm-api/helpers"
	"github.com/unnamedxaer/gymm-api/usecases"
	"github.com/unnamedxaer/gymm-api/validation"
)

func trimWhitespacesOnExerciseInput(e *usecases.ExerciseInput) {
	e.Name = helpers.TrimWhiteSpaces(e.Name)
	e.Description = helpers.TrimWhiteSpaces(e.Description)
}

func trimWhitespacesOnUserInput(u *usecases.UserInput) {
	u.Username = helpers.TrimWhiteSpaces(u.Username)
}

func getExerciseExcludeTags(e *usecases.ExerciseInput) []string {
	excludedTags := make([]string, 2)
	tag, _ := validation.GetFieldJSONTag(e, "CreatedBy")
	excludedTags[0] = tag
	tag, _ = validation.GetFieldJSONTag(e, "CreatedAt")
	excludedTags[1] = tag
	return excludedTags
}
