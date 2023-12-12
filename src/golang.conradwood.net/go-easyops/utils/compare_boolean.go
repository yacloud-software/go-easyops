package utils

import (
	"fmt"
	"strings"
)

// return true or false if a boolean can be parsed. error if string is not something that is understood as a boolean. currently understands true/false,yes/no and on/off
func BooleanValue(boolvalue string) (bool, error) {
	s := strings.ToLower(boolvalue)
	if s == "true" || s == "yes" || s == "on" {
		return true, nil
	}
	if s == "false" || s == "no" || s == "off" {
		return false, nil
	}
	return false, fmt.Errorf("string \"%s\" is not a boolean value", boolvalue)
}

// if a valid boolean can parsed, return true, false in all other cases
func BooleanValueNoErr(boolvalue string) bool {
	s := strings.ToLower(boolvalue)
	if s == "true" || s == "yes" || s == "on" {
		return true
	}
	return false
}
