// Package requests contains general utility functions for tournament conversion
package requests

// In JavaScript: Math.trunc(Math.log2(n))
// Do not use with negative numbers.
// See: https://stackoverflow.com/questions/19339594/truncated-binary-logarithm
func TruncLog2(n int) (exp, value int) {
	if n <= 0 {
		return
	}
	value = 1
	n = n >> 1
	for n != 0 {
		exp++
		value *= 2
		n = n >> 1
	}
	return
}

func CalculatePlacings(numEntrants int) (placings int) {
	// Annoying edge cases
	if numEntrants < 1 {
		return -1
	}
	if numEntrants <= 4 {
		return numEntrants
	}

	exp, value := TruncLog2(numEntrants)
	placings = 2*exp
	if numEntrants <= (3*value)/2 {
		placings = placings + 1
	}
	if numEntrants > (3*value)/2 {
		placings = placings + 2
	}
	return
}