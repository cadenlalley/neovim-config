package override

import "gopkg.in/guregu/null.v4"

func String(input string, original string) string {
	if input != "" {
		return input
	}
	return original
}

func NullString(input string, original null.String) (output null.String) {
	if input != original.String {
		return null.NewString(input, input != "")
	}
	return original
}

func Bool(input *bool, original bool) bool {
	if input != nil && *input != original {
		return *input
	}
	return original
}
