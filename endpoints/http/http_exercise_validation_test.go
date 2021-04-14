package http

import (
	"fmt"
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
		// {
		// 	desc: "valid exercise input",
		// 	input: &usecases.ExerciseInput{
		// 		Name:        mocks.ExampleExercise.Name,
		// 		Description: mocks.ExampleExercise.Description,
		// 		SetUnit:     mocks.ExampleExercise.SetUnit,
		// 	},
		// 	want: nil,
		// },
		// {
		// 	desc: "valid exercise input 2",
		// 	input: &usecases.ExerciseInput{
		// 		Name:        mocks.ExampleExercise.Name + "2",
		// 		Description: mocks.ExampleExercise.Description,
		// 		SetUnit:     mocks.ExampleExercise.SetUnit,
		// 	},
		// 	want: nil,
		// },
		// {
		// 	desc: "valid exercise, two words name",
		// 	input: &usecases.ExerciseInput{
		// 		Name:        "Front Squat",
		// 		Description: mocks.ExampleExercise.Description,
		// 		SetUnit:     mocks.ExampleExercise.SetUnit,
		// 	},
		// 	want: nil,
		// },
		// {
		// 	desc: "too long name",
		// 	input: &usecases.ExerciseInput{
		// 		Name:        "too long name too long name too long name too long name",
		// 		Description: mocks.ExampleExercise.Description,
		// 		SetUnit:     mocks.ExampleExercise.SetUnit,
		// 	},
		// 	want: []string{"'name'"},
		// },
		// {
		// 	desc: "too short name",
		// 	input: &usecases.ExerciseInput{
		// 		Name:        "t",
		// 		Description: mocks.ExampleExercise.Description,
		// 		SetUnit:     mocks.ExampleExercise.SetUnit,
		// 	},
		// 	want: []string{"'name'"},
		// },
		// {
		// 	desc: "name contains forbidden characters",
		// 	input: &usecases.ExerciseInput{
		// 		Name:        "DL ╤",
		// 		Description: mocks.ExampleExercise.Description,
		// 		SetUnit:     mocks.ExampleExercise.SetUnit,
		// 	},
		// 	want: []string{"'name'"},
		// },
		// {
		// 	desc: "missing name",
		// 	input: &usecases.ExerciseInput{
		// 		Description: mocks.ExampleExercise.Description,
		// 		SetUnit:     mocks.ExampleExercise.SetUnit,
		// 	},
		// 	want: []string{"'name'"},
		// },
		{
			desc: "not printable chars in description: 💤",
			input: &usecases.ExerciseInput{
				Name:        mocks.ExampleExercise.Name,
				Description: mocks.ExampleExercise.Description + "💤",
				SetUnit:     mocks.ExampleExercise.SetUnit,
			},
			want: []string{"'description'"},
		},
		{
			desc: "not printable chars in description 2: ╨",
			input: &usecases.ExerciseInput{
				Name:        mocks.ExampleExercise.Name,
				Description: mocks.ExampleExercise.Description + "╨",
				SetUnit:     mocks.ExampleExercise.SetUnit,
			},
			want: []string{"'description'"},
		},
		{
			desc: "not printable chars in description 3: \t(tab)",
			input: &usecases.ExerciseInput{
				Name:        mocks.ExampleExercise.Name,
				Description: mocks.ExampleExercise.Description + "\t 2223",
				SetUnit:     mocks.ExampleExercise.SetUnit,
			},
			want: []string{"'description'"},
		},
		{
			desc: "not printable chars in description 4: t",
			input: &usecases.ExerciseInput{
				Name:        mocks.ExampleExercise.Name,
				Description: (mocks.ExampleExercise.Description + "   t"),
				SetUnit:     mocks.ExampleExercise.SetUnit,
			},
			want: nil,
		},
		// {
		// 	desc: "missing description and set unit",
		// 	input: &usecases.ExerciseInput{
		// 		Name: mocks.ExampleExercise.Name,
		// 	},
		// 	want: []string{"'description'", "'setUnit'"},
		// },
		// {
		// 	desc: "missing name, description and incorrect set unit",
		// 	input: &usecases.ExerciseInput{
		// 		SetUnit: 12,
		// 	},
		// 	want: []string{"'description'", "'setUnit'", "'name'"},
		// },
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

			errText := err.Error()
			for _, w := range tC.want {
				if !strings.Contains(errText, w) {
					if err == nil {
						fmt.Println("want not nil error")
					}
					t.Errorf("want %q, got %v, for %+v", w, err, tC.input)
				}
			}
		})
	}
}
