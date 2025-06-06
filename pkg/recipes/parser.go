package recipes

import (
	"strings"

	"gopkg.in/guregu/null.v4"
)

// var ValidUnits map[string][]string = map[string][]string{
// 	"bag":        {},
// 	"bottle":     {},
// 	"box":        {},
// 	"can":        {},
// 	"clove":      {},
// 	"cup":        {"c"},
// 	"dash":       {},
// 	"drop":       {},
// 	"gallon":     {},
// 	"gram":       {},
// 	"jar":        {},
// 	"kilogram":   {},
// 	"liter":      {},
// 	"milliliter": {},
// 	"ounce":      {"oz"},
// 	"packet":     {},
// 	"piece":      {},
// 	"pint":       {},
// 	"pinch":      {},
// 	"pound":      {"#", "lb"},
// 	"quart":      {},
// 	"slice":      {},
// 	"stick":      {},
// 	"tbsp":       {"T"},
// 	"tsp":        {},
// }

// func getValidUnit(input string) string {
// 	for unit, aliases := range ValidUnits {
// 		if len(aliases) > 0 {
// 			for _, alias := range aliases {
// 				if alias == input {
// 					return unit
// 				}
// 			}
// 		}
// 	}

// 	return input
// }

// func ParseIngredientUnit(input null.String) null.String {
// 	unit := ParseNullString(input)
// 	if unit.IsZero() {
// 		return unit
// 	}

// 	output := getValidUnit(unit.String)

// 	return null.NewString(output, output != "")
// }

func ParseNullString(input null.String) null.String {
	if input.IsZero() {
		return input
	}

	invalids := []string{"null", "-", "", "n/a"}
	for _, v := range invalids {
		if strings.ToLower(input.String) == v {
			return null.String{}
		}
	}

	return input
}

func ParseNullFloat(input null.Float) null.Float {
	if input.IsZero() {
		return input
	}

	if input.Float64 == 0 {
		return null.Float{}
	}

	return input
}
