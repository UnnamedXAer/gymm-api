package helpers

import "testing"

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
		got := StrSliceIndexOf(slice, given)
		if got != wanted {
			t.Errorf("Expected to get %d, got %d, str: %q", got, wanted, given)
		}
	}

}
