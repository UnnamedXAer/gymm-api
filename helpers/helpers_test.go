package helpers_test

import (
	"testing"

	"github.com/unnamedxaer/gymm-api/helpers"
)

var (
	slice []string = []string{"test", "development", "production"}
)

func TestStrSliceIndexOf(t *testing.T) {
	givenWanted := map[string]int{
		slice[0]:            0,
		slice[len(slice)-1]: len(slice) - 1,
		"":                  -1,
		slice[0] + "X":      -1,
	}

	for given, wanted := range givenWanted {
		got := helpers.StrSliceIndexOf(slice, given)
		if got != wanted {
			t.Errorf("Expected to get %d, got %d, str: %q", got, wanted, given)
		}
	}

}

func TestTrimWhiteSpacesJoinFields(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  string
	}{
		{
			desc:  "",
			input: "asd",
			want:  "asd",
		},
		{
			desc:  "",
			input: "asd ",
			want:  "asd",
		},
		{
			desc:  "",
			input: "as d ",
			want:  "as d",
		},
		{
			desc:  "",
			input: " as d ",
			want:  "as d",
		},
		{
			desc:  "",
			input: "\n as d ",
			want:  "as d",
		},
		{
			desc:  "",
			input: "\n as d\t",
			want:  "as d",
		},
		{
			desc:  "",
			input: "as\td",
			want:  "as d",
		},
		{
			desc:  "",
			input: "as\rd",
			want:  "as d",
		},
		{
			desc:  "",
			input: "asd 1",
			want:  "asd 1",
		},
		{
			desc:  "",
			input: "asd  1",
			want:  "asd 1",
		},
		{
			desc:  "",
			input: "\tThe Grow method can be  used to preallocate memory\n when the maximum size of the\r string is known.",
			want:  "The Grow method can be used to preallocate memory when the maximum size of the string is known.",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := helpers.TrimWhiteSpaces(tC.input)
			if got != tC.want {
				t.Errorf("want %q, got %q", tC.want, got)
			}
		})
	}
}

func BenchmarkTrimWhiteSpacesJoinFields(b *testing.B) {
	s := "\tThe Grow method can be used to preallocate memory\n when the maximum size of the\r string is known."
	for i := 0; i < b.N; i++ {
		helpers.TrimWhiteSpaces(s)
	}
}
