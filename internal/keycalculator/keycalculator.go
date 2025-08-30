package keycalculator

import "fmt"

// KeyCalculator returns the number of keys associated with a given key length.
type KeyCalculator func(keyLength uint) (int, error)

// DefaultCalculate returns the expected number of keys given known key lengths.
func DefaultCalculate(keyLength uint) (int, error) {
	keyLengths := map[uint]int{
		24: 16777216,   // 2^24
		28: 268435456,  // 2^28
		32: 4294967296, // 2^32
	}
	
	if numKeys, ok := keyLengths[keyLength]; ok {
		return numKeys, nil
	}
	return 0, fmt.Errorf("unexpected key length %d: valid key lengths are 24, 28 and 32 bits", keyLength)
}
