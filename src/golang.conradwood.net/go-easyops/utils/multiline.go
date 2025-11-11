package utils

import (
	"slices"
	"strings"
)

/*
Adds a prefix to each line of a multi-line string.
it hounours existing line breaks (\n)
if line is longer than linelen it will add a linebreak
empty trailing lines will be trimmed
a single trailing \n will be retained.
*/
func MultiLinePrefix(input, prefix string, linelen int) string {
	lines := strings.Split(input, "\n")
	if len(lines) == 0 {
		return ""
	}

	// remove trailing empty lines
	for lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
		if len(lines) == 0 {
			return ""
		}
	}

	// word wrap
	repeat := true
	for repeat {
		repeat = false
		for i, line := range lines {
			if len(line) <= linelen {
				continue
			}
			oline := line[:linelen]
			rline := line[linelen:]

			lines[i] = rline
			lines = slices.Insert(lines, i, oline)
			repeat = true
		}
	}

	for i, line := range lines {
		lines[i] = prefix + line
	}
	res := strings.Join(lines, "\n")
	if input[len(input)-1] == '\n' {
		res = res + "\n"
	}
	return res
}
