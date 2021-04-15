package http

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/unnamedxaer/gymm-api/mocks"
	"github.com/unnamedxaer/gymm-api/usecases"
	"github.com/unnamedxaer/gymm-api/validation"
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
			desc: "valid exercise input 2",
			input: &usecases.ExerciseInput{
				Name:        mocks.ExampleExercise.Name + "2",
				Description: mocks.ExampleExercise.Description,
				SetUnit:     mocks.ExampleExercise.SetUnit,
			},
			want: nil,
		},
		{
			desc: "valid exercise, two words name",
			input: &usecases.ExerciseInput{
				Name:        "Front Squat",
				Description: mocks.ExampleExercise.Description,
				SetUnit:     mocks.ExampleExercise.SetUnit,
			},
			want: nil,
		},
		{
			desc: "too long name",
			input: &usecases.ExerciseInput{
				Name:        "too long name too long name too long name too long name",
				Description: mocks.ExampleExercise.Description,
				SetUnit:     mocks.ExampleExercise.SetUnit,
			},
			want: []string{"name"},
		},
		{
			desc: "too short name",
			input: &usecases.ExerciseInput{
				Name:        "t",
				Description: mocks.ExampleExercise.Description,
				SetUnit:     mocks.ExampleExercise.SetUnit,
			},
			want: []string{"name"},
		},
		{
			desc: "name contains forbidden characters",
			input: &usecases.ExerciseInput{
				Name:        "DL ╤",
				Description: mocks.ExampleExercise.Description,
				SetUnit:     mocks.ExampleExercise.SetUnit,
			},
			want: []string{"name"},
		},
		{
			desc: "missing name",
			input: &usecases.ExerciseInput{
				Description: mocks.ExampleExercise.Description,
				SetUnit:     mocks.ExampleExercise.SetUnit,
			},
			want: []string{"name"},
		},
		{
			desc: "not printable chars in description: 💤",
			input: &usecases.ExerciseInput{
				Name:        mocks.ExampleExercise.Name,
				Description: mocks.ExampleExercise.Description + "💤",
				SetUnit:     mocks.ExampleExercise.SetUnit,
			},
			want: []string{"description"},
		},
		{
			desc: "not printable chars in description 2: ╨",
			input: &usecases.ExerciseInput{
				Name:        mocks.ExampleExercise.Name,
				Description: mocks.ExampleExercise.Description + "╨",
				SetUnit:     mocks.ExampleExercise.SetUnit,
			},
			want: []string{"description"},
		},
		{
			desc: "not printable chars in description 3: \t(tab)",
			input: &usecases.ExerciseInput{
				Name:        mocks.ExampleExercise.Name,
				Description: mocks.ExampleExercise.Description + "\t 2223",
				SetUnit:     mocks.ExampleExercise.SetUnit,
			},
			want: []string{"description"},
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
		{
			desc: "missing description and set unit",
			input: &usecases.ExerciseInput{
				Name: mocks.ExampleExercise.Name,
			},
			want: []string{"description", "setUnit"},
		},
		{
			desc: "missing name, description and incorrect set unit",
			input: &usecases.ExerciseInput{
				SetUnit: 12,
			},
			want: []string{"description", "setUnit", "name"},
		},
	}
	runExerciseTestCases(t, validateExerciseInput, testCases)
}

func TestValidateExerciseInput4Update(t *testing.T) {
	testCases := []struct {
		desc  string
		input *usecases.ExerciseInput
		want  []string
	}{
		{
			desc: "valid, only name",
			input: &usecases.ExerciseInput{
				Name: mocks.ExampleExercise.Name,
			},
			want: nil,
		},
		{
			desc: "incorrect name, description, set unit",
			input: &usecases.ExerciseInput{
				Name:        mocks.ExampleExercise.Name + " 🛑 ",
				Description: mocks.ExampleExercise.Description + " 🛑 ",
				SetUnit:     mocks.ExampleExercise.SetUnit + 12,
			},
			want: []string{"name", "description", "setUnit"},
		},
		{
			desc: "incorrect name, valid set unit",
			input: &usecases.ExerciseInput{
				Name:    "X",
				SetUnit: mocks.ExampleExercise.SetUnit,
			},
			want: []string{"name"},
		},
		{
			desc: "valid, only description",
			input: &usecases.ExerciseInput{
				Description: "The Overhead Press is a full body, compound exercise.",
			},
			want: nil,
		},
	}
	runExerciseTestCases(t, validateExerciseInput4Update, testCases)
}

func runExerciseTestCases(
	t *testing.T,
	testedFunc func(validate *validator.Validate, exercise *usecases.ExerciseInput) error,
	testCases []struct {
		desc  string
		input *usecases.ExerciseInput
		want  []string
	}) {
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := testedFunc(validate, tC.input)

			if tC.want == nil {
				if err != nil {
					t.Errorf("want nil error, got %v, for %+v", err, tC.input)
				}
				return
			}

			strValErr, ok := err.(*validation.StructValidError)
			if !ok {
				t.Errorf("want errors for fields %v, got %v, for %+v", tC.want, err, tC.input)
				return
			}

			vErrs := strValErr.ValidationErrors()
			if len(vErrs) != len(tC.want) {
				t.Errorf("want errors (%d) for fields %v, got (%d) %v, for %+v",
					len(tC.want), tC.want, len(vErrs), vErrs, tC.input)
				return
			}

			var failed bool
			for _, f := range tC.want {
				_, ok = vErrs[f]
				if !ok {
					t.Errorf("missing error for %q", f)
					failed = true
				}
			}

			if failed {
				t.Logf("errors: %+v", vErrs)
				t.Logf("input: %+v", tC.input)
			}
		})
	}
}
