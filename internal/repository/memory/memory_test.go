package memory_test

import (
	"testing"

	"github.com/andream16/mitmcracker/internal/repository"
	"github.com/andream16/mitmcracker/internal/repository/memory"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemo_Insert(t *testing.T) {
	t.Run("it should successfully find a common ciphertext", func(t *testing.T) {
		inMemo := &memory.InMemo{
			KeyPairs: map[string]*repository.KeyPair{},
		}

		const (
			matchingEncKey     = "matchingEncKey"
			matchingDecKey     = "matchingDecKey"
			matchingCipherText = "matchingCiphertext"
		)

		for _, in := range []struct {
			key           string
			mode          string
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

func TestInMemo_InsertBulk(t *testing.T) {
	t.Run("it should successfully find a common ciphertext", func(t *testing.T) {
		inMemo := &memory.InMemo{
			KeyPairs: map[string]*repository.KeyPair{},
		}

		const (
			matchingEncKey     = "matchingEncKey"
			matchingDecKey     = "matchingDecKey"
			matchingCipherText = "matchingCiphertext"
		)

		insertReqs := []repository.InsertBulkRequest{
			{
				Key:        "notMatchingEncKey1",
				Mode:       repository.EncodeMode,
				CipherText: "notMatchingCipherText1",
			},
			{
				Key:        "notMatchingEncKey2",
				Mode:       repository.EncodeMode,
				CipherText: "notMatchingCipherText2",
			},
			{
				Key:        matchingDecKey,
				Mode:       repository.DecodeMode,
				CipherText: matchingCipherText,
			},
			{
				Key:        "NotMatchingDecKey4",
				Mode:       repository.DecodeMode,
				CipherText: "notMatchingCipherText4",
			},
			{
				Key:        matchingEncKey,
				Mode:       repository.EncodeMode,
				CipherText: matchingCipherText,
			},
		}

		kp, found, err := inMemo.InsertBulk(insertReqs...)
		require.NoError(t, err)
		require.True(t, found)
		assert.Equal(t, matchingEncKey, kp.EncodeKey)
		assert.Equal(t, matchingDecKey, kp.DecodeKey)
	})
}

// BenchmarkInMemo_Insert-8   	18032630	        79.5 ns/op
func BenchmarkInMemo_Insert(b *testing.B) {
	inMemo := &memory.InMemo{
		KeyPairs: map[string]*repository.KeyPair{},
	}

	const (
		matchingEncKey     = "matchingEncKey"
		matchingDecKey     = "matchingDecKey"
		matchingCipherText = "matchingCiphertext"
	)

	insertReqs := []struct {
		key        string
		mode       string
		cipherText string
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
			key:        matchingEncKey,
			mode:       repository.EncodeMode,
			cipherText: matchingCipherText,
		},
	}

	for n := 0; n < b.N; n++ {
		for _, r := range insertReqs {
			_, _, _ = inMemo.Insert(r.key, r.cipherText, r.mode)
		}
	}
}

// BenchmarkInMemo_InsertBulk-8   	99254767	        11.2 ns/op
func BenchmarkInMemo_InsertBulk(b *testing.B) {
	inMemo := &memory.InMemo{
		KeyPairs: map[string]*repository.KeyPair{},
	}

	const (
		matchingEncKey     = "matchingEncKey"
		matchingDecKey     = "matchingDecKey"
		matchingCipherText = "matchingCiphertext"
	)

	insertReqs := []repository.InsertBulkRequest{
		{
			Key:        "notMatchingEncKey1",
			Mode:       repository.EncodeMode,
			CipherText: "notMatchingCipherText1",
		},
		{
			Key:        "notMatchingEncKey2",
			Mode:       repository.EncodeMode,
			CipherText: "notMatchingCipherText2",
		},
		{
			Key:        matchingDecKey,
			Mode:       repository.DecodeMode,
			CipherText: matchingCipherText,
		},
		{
			Key:        "NotMatchingDecKey4",
			Mode:       repository.DecodeMode,
			CipherText: "notMatchingCipherText4",
		},
		{
			Key:        matchingEncKey,
			Mode:       repository.EncodeMode,
			CipherText: matchingCipherText,
		},
	}

	for n := 0; n < b.N; n++ {
		_, _, _ = inMemo.InsertBulk(insertReqs...)
	}
}
