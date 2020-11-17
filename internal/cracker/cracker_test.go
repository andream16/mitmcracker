package cracker_test

import (
	"context"
	"testing"

	"github.com/andream16/mitmcracker/internal/repository"

	"github.com/andream16/mitmcracker/internal/cracker"
	"github.com/andream16/mitmcracker/internal/perf"
	"github.com/andream16/mitmcracker/internal/repository/memory"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCracker_Crack(t *testing.T) {
	const (
		keyLength        uint = 24
		knownPlainText        = "C330C9CBD01DFBA0"
		knownEncodedText      = "E10C65124518DB05"
		keysNumber            = 10
		commonCypherText      = "hello"
		matchingKey           = "000006"
	)

	var (
		decFn = func(key, plainText string) (string, error) {
			k, _ := map[string]string{
				"000000":    "AAAAAA",
				"000001":    "AAAAAB",
				"000002":    "AAAAAC",
				"000003":    "AAAAAD",
				"000004":    "AAAAAE",
				"000005":    "AAAAAF",
				matchingKey: commonCypherText,
				"000007":    "AAAABB",
				"000008":    "AAAABC",
				"000009":    "AAAABD",
			}[key]
			return k, nil
		}
		encFn = func(key, plainText string) (string, error) {
			k, _ := map[string]string{
				"000000":    "FFFFFF",
				"000001":    "FFFFFE",
				"000002":    "FFFFFD",
				"000003":    "FFFFFC",
				"000004":    "FFFFFB",
				"000005":    "FFFFFA",
				matchingKey: commonCypherText,
				"000007":    "FFFFEB",
				"000008":    "FFFFDC",
				"000009":    "FFFFDE",
			}[key]
			return k, nil
		}
		formatterFn = func(key int, keyLength uint) string {
			k, _ := map[int]string{
				0: "000000",
				1: "000001",
				2: "000002",
				3: "000003",
				4: "000004",
				5: "000005",
				6: matchingKey,
				7: "000007",
				8: "000008",
				9: "000009",
			}[key]
			return k
		}
		keyCalcFn = func(keyLength uint) (int, error) {
			return keysNumber, nil
		}
		maxGoRoutinesNumFn = perf.DefaultMaxGoRoutineNumber
	)

	cr, err := cracker.New(
		keyLength,
		knownEncodedText,
		knownPlainText,
		&memory.InMemo{
			KeyPairs: map[string]*repository.KeyPair{},
		},
		encFn,
		decFn,
		formatterFn,
		keyCalcFn,
		maxGoRoutinesNumFn,
		2,
	)
	require.NoError(t, err)

	kp, found, err := cr.Crack(context.Background())
	require.NoError(t, err)
	assert.True(t, found)
	require.NotNil(t, kp)

	assert.Equal(t, matchingKey, kp.DecodeKey, kp.EncodeKey)
}
