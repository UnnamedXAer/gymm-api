package http

import (
	"strings"
	"testing"

	"github.com/unnamedxaer/gymm-api/mocks"
	"github.com/unnamedxaer/gymm-api/usecases"
)

func TestValidateExerciseInput(t *testing.T) {
	testCases := []struct {
		desc  string
		input *usecases.ExerciseInput
		want  []string
	}{
		{
			desc: "valid exercise input",
			input: &usecases.ExerciseInput{
				Name:        mocks.ExampleExercise.Name,
				Description: mocks.ExampleExercise.Description,
				SetUnit:     mocks.ExampleExercise.SetUnit,
			},
			want: nil,
		},
		{
			desc: "missing name",
			input: &usecases.ExerciseInput{
				Description: mocks.ExampleExercise.Description,
				SetUnit:     mocks.ExampleExercise.SetUnit,
			},
			want: []string{"'name'"},
		},
		{
			desc: "missing description and set unit",
			input: &usecases.ExerciseInput{
				Name: mocks.ExampleExercise.Name,
			},
			want: []string{"'description'", "'set_unit'"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := validateExerciseInput(validate, tC.input)
			if tC.want == nil {
				if err != nil {
					t.Errorf("want nil error, got %v, for %+v", err, tC.input)
				}
				return
			}
			for _, w := range tC.want {
				if !strings.Contains(err.Error(), w) {
					t.Errorf("want %q, got %v, for %+v", w, err, tC.input)
				}
			}
		})
	}
}
