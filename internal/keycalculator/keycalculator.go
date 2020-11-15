package keycalculator

import "fmt"

// KeyCalculator returns the number of keys associated with a given key length.
type KeyCalculator func(keyLength uint) (int, error)

// DefaultCalculate returns the expected number of keys given known key lengths.
func DefaultCalculate(keyLength uint) (int, error) {
	if l, ok := map[uint]int{
		24: 16777216,
		28: 268435456,
		32: 4294967296,
	}[keyLength]; ok {
		return l, nil
	}
	return 0, fmt.Errorf("unexpected key length %d. Valid key legths are 24, 28 and 32 bits", keyLength)
}
