package conversionutil

import (
	"strconv"
	"strings"
)

/*
Uint8toString converts a uint8 array to a string array

ex.
[91 34 97 34 44 32 34 98 34 44 32 34 99 34 93] -> ["a", "b", "c"]
*/
func Uint8toString(array []uint8) []string {
	stringNoBrackets := string(array[1 : len(array)-1])
	return strings.Split(stringNoBrackets, ",")
}

/*
RemoveDuplicateValues removes duplicate values from a string array

ex.
["a", "b", "c", "a", "b"] -> ["a", "b", "c"]
*/
func RemoveDuplicateValues(values []string) (cleaned []string) {
	keys := make(map[string]struct{})
	for _, entry := range values {
		if _, value := keys[entry]; !value {
			keys[entry] = struct{}{}
			cleaned = append(cleaned, entry)
		}
	}
	return cleaned
}

/*
ConvertLargeStringNum is a nifty lil conversion function that parses string ints that would overwhelm int64 mem space, handling e notation and decimals, or lack thereof

ex.
"1.234567890123456789012345678901234567890123456789012345678901234567890123456789e+100" -> "123456789012345678901234567890123456789012345678901234567890123456789012345678900000000000
*/
func ConvertLargeStringNum(txt string) (converted string) {
	ePosition := strings.Index(txt, "e")
	// there is no 'e' in the string
	if ePosition == -1 {
		return txt
	}
	distance, _ := strconv.Atoi(txt[ePosition+2:])
	decimalPosition := strings.Index(txt, ".")
	// no decimal, just pad with zeroes
	if decimalPosition == -1 {
		return convertNoDecimal(txt[:ePosition], distance)
	}
	return moveDecimal(txt[:decimalPosition], txt[decimalPosition+1:ePosition], string(txt[ePosition+1]), distance)
}

/*
convertNoDecimal pads the string with zeroes to the right of the string to the distance specified

ex.
"12345678901234567890123456789012345678901234567890123456789012345678901234567890" -> "123456789012345678901234567890123456789012345678901234567890123456789012345678900000000000
*/
func convertNoDecimal(txt string, distance int) (converted string) {
	converted = txt
	for i := 0; i < distance; i++ {
		converted += "0"
	}
	return converted
}

/*
moveDecimal moves the decimal point in the string to the right or left by the distance specified
*/
func moveDecimal(txtBeforeDecimal, txtAfterDecimal, direction string, distance int) (converted string) {
	converted = txtBeforeDecimal
	if len(txtAfterDecimal) == distance {
		converted += txtAfterDecimal
	}
	if len(txtAfterDecimal) < distance {
		converted += txtAfterDecimal
		for i := 0; i < distance-len(txtAfterDecimal); i++ {
			converted += "0"
		}
	}
	return converted
}
