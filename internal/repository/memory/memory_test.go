package memory_test

import (
	"testing"

	"github.com/andream16/mitmcracker/internal/repository"
	"github.com/andream16/mitmcracker/internal/repository/memory"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemo_Insert(t *testing.T) {
	t.Run("it should return an error because the mode is not valid", func(t *testing.T) {
		inMemo := &memory.InMemo{}
		kp, found, err := inMemo.Insert("", "someMode", "")
		require.Error(t, err)
		assert.False(t, found)
		assert.Nil(t, kp)
	})
	t.Run("it should successfully find a common ciphertext", func(t *testing.T) {
		inMemo := &memory.InMemo{}

		const (
			matchingEncKey     = "matchingEncKey"
			matchingDecKey     = "matchingDecKey"
			matchingCipherText = "matchingCiphertext"
		)

		for _, in := range []struct {
			key           string
			mode          repository.Mode
			cipherText    string
			expectedFound bool
		}{
			{
				key:        "notMatchingEncKey1",
				mode:       repository.EncodeMode,
				cipherText: "notMatchingCipherText1",
			},
			{
				key:        "notMatchingEncKey2",
				mode:       repository.EncodeMode,
				cipherText: "notMatchingCipherText2",
			},
			{
				key:        matchingDecKey,
				mode:       repository.DecodeMode,
				cipherText: matchingCipherText,
			},
			{
				key:        "NotMatchingDecKey4",
				mode:       repository.DecodeMode,
				cipherText: "notMatchingCipherText4",
			},
			{
				key:           matchingEncKey,
				mode:          repository.EncodeMode,
				cipherText:    matchingCipherText,
				expectedFound: true,
			},
		} {
			kp, found, err := inMemo.Insert(in.key, in.cipherText, in.mode)
			require.NoError(t, err)

			if in.expectedFound {
				require.True(t, found)
				assert.Equal(t, matchingEncKey, kp.EncodeKey)
				assert.Equal(t, matchingDecKey, kp.DecodeKey)
				continue
			}

			assert.False(t, found)
			assert.Nil(t, kp)
		}
	})
}
