// Package implements an internal mechanism to communicate with an impress terminal.
package fontspec

// Values for "style" atribute
var styleValues = map[string]int{
	"":        0,
	"normal":  0,
	"oblique": 1,
	"italic":  2,
}

// Values for "variant" atribute
var variantValues = map[string]int{
	"":           0,
	"normal":     0,
	"small_caps": 1,
	"small caps": 1,
}

// Values for "weight" atribute
var weightValues = map[string]int{
	"thin":       100,
	"ultralight": 200,
	"light":      300,
	"semilight":  350,
	"book":       380,
	"":           400,
	"normal":     400,
	"medium":     500,
	"semibold":   600,
	"bold":       700,
	"ultrabold":  800,
	"heavy":      900,
	"ultraheavy": 1000,
}

// Values for "stretch" atribute
var stretchValues = map[string]int{
	"ultra_condensed": 0,
	"ultra condensed": 0,
	"extra_condensed": 1,
	"extra condensed": 1,
	"condensed":       2,
	"semi_condensed":  3,
	"semi condensed":  3,
	"":                4,
	"normal":          4,
	"semi_expanded":   5,
	"semi expanded":   5,
	"expanded":        6,
	"extra_expanded":  7,
	"extra expanded":  7,
	"ultra_expanded":  8,
	"ultra expanded":  8,
}

func Attributes(attributes map[string]string) (string, int, int, int, int) {
	family := attributes["family"]
	style := 0
	if value, ok := styleValues[attributes["style"]]; ok {
		style = value
	}
	variant := 0
	if value, ok := variantValues[attributes["variant"]]; ok {
		variant = value
	}
	weight := 400
	if value, ok := weightValues[attributes["weight"]]; ok {
		weight = value
	}
	stretch := 4
	if value, ok := stretchValues[attributes["stretch"]]; ok {
		stretch = value
	}
	return family, style, variant, weight, stretch
}

func SplitByLengths(text string, lengths []int) []string {
	output := make([]string, 0, len(lengths))
	for _, length := range lengths {
		if length > len(text) {
			break
		}
		output = append(output, text[:length])
		text = text[length:]
	}
	return output
}
