// Package implements an internal mechanism to communicate with an impress terminal.
package fontspec

// Values for "style" atribute
func styleValues(style string) int {
	switch style {
	case "normal":
		return 0
	case "oblique":
		return 1
	case "italic":
		return 2
	default:
		return 0
	}
}

// Values for "variant" atribute
func variantValues(variant string) int {
	switch variant {
	case "normal":
		return 0
	case "small_caps", "small caps":
		return 1
	default:
		return 0
	}
}

// Values for "weight" atribute
func weightValues(weight string) int {
	switch weight {
	case "thin":
		return 100
	case "ultralight":
		return 200
	case "light":
		return 300
	case "semilight":
		return 350
	case "book":
		return 380
	case "normal":
		return 400
	case "medium":
		return 500
	case "semibold":
		return 600
	case "bold":
		return 700
	case "ultrabold":
		return 800
	case "heavy":
		return 900
	case "ultraheavy":
		return 1000
	default:
		return 400
	}
}

// Values for "stretch" atribute
func stretchValues(stretch string) int {
	switch stretch {
	case "ultra_condensed", "ultra condensed":
		return 0
	case "extra_condensed", "extra condensed":
		return 1
	case "condensed":
		return 2
	case "semi_condensed", "semi condensed":
		return 3
	case "normal":
		return 4
	case "semi_expanded", "semi expanded":
		return 5
	case "expanded":
		return 6
	case "extra_expanded", "extra expanded":
		return 7
	case "ultra_expanded", "ultra expanded":
		return 8
	default:
		return 4
	}
}

func Attributes(attributes map[string]string) (string, int, int, int, int) {
	return attributes["family"],
		styleValues(attributes["style"]),
		variantValues(attributes["variant"]),
		weightValues(attributes["weight"]),
		stretchValues(attributes["stretch"])
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
