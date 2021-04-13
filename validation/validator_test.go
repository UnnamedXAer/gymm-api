package validation

import (
	"reflect"
	"testing"
)

func TestValidatePassword(t *testing.T) {

	givenWanted := map[string]bool{
		"":                                     false,
		"1":                                    false,
		"1a":                                   false,
		"111111":                               false,
		"123456":                               false,
		"1234567":                              false,
		"123456a":                              true,
		"aaaaaa":                               false,
		"asdfghj":                              false,
		"asdfghj1":                             true,
		"AaSsDdFf":                             false,
		"AaSsDdFf1":                            true,
		"!!!!!!!":                              false,
		"asdfg":                                false,
		"asdfg^":                               true,
		"112233aa":                             false,
		"11223344aa":                           true,
		"AAAAAAAAAsssssssssssssssssssAAAAAA12": false,
		"AAAAAAAAAsssssssssssssssssssAAAAAAAAAAAAA12":  false,
		"AAAAAAAAAsssssssssssssssssssSAAAAAAAAAAAAA12": true,
		"!@##$$%^%":            false,
		"!@##2$%^%":            true,
		"!@##2a%^%":            true,
		"Q!W@E#R$":             true,
		"12444%":               false,
		"12344%":               true,
		"asdffø":               true,
		"!!!!@@@@####$$$$%%%%": false,
		"!!!!@@@@####$$$$øøøø": true,
	}

	for given, wanted := range givenWanted {
		got := validatePassword(given)
		if got != wanted {

			got = validatePassword(given)

		}
		if got != wanted {
			t.Errorf("password: '%s', want: %t, got: %t", given, wanted, got)
		}
	}
}

func TestValidateSetUnit(t *testing.T) {

	givenWanted := map[interface{}]bool{
		-1: false,
		0:  false,
		1:  true,
		2:  true,
		3:  false,
	}

	for input, want := range givenWanted {
		got := validateSetUnit(reflect.ValueOf(input))
		if got != want {
			t.Errorf("set unit: %v, want: %t, got: %t", input, want, got)
		}
	}
}
