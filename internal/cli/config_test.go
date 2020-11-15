package cli_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/andream16/mitmcracker/internal/cli"

	"github.com/stretchr/testify/require"
)

func TestInput_validate(t *testing.T) {
	t.Run("should fail because key length is not supported", func(t *testing.T) {
		conf := &cli.Config{}
		err := conf.Validate()
		require.Error(t, err)

		var e cli.InvalidInputFlagError
		require.True(t, errors.As(err, &e))
	})
	t.Run("should fail because key length is invalig", func(t *testing.T) {
		conf := &cli.Config{
			KeyLength: 10,
		}
		err := conf.Validate()
		require.Error(t, err)

		var e cli.InvalidInputFlagError
		require.True(t, errors.As(err, &e))
	})
	t.Run("should fail because known encrypted text is empty", func(t *testing.T) {
		conf := &cli.Config{
			KeyLength: 24,
		}
		err := conf.Validate()
		require.Error(t, err)

		var e cli.InvalidInputFlagError
		require.True(t, errors.As(err, &e))
	})
	t.Run("should fail because known plain text is empty", func(t *testing.T) {
		conf := &cli.Config{
			KeyLength: 28,
			EncText:   "someEncryptedText",
		}
		err := conf.Validate()
		require.Error(t, err)

		var e cli.InvalidInputFlagError
		require.True(t, errors.As(err, &e))
	})
	t.Run("should fail because storage is not supported", func(t *testing.T) {
		conf := &cli.Config{
			KeyLength: 28,
			EncText:   "someEncryptedText",
			PlainText: "somePlainText",
			Storage: &cli.Storage{
				Type: "xyz",
			},
		}
		err := conf.Validate()
		require.Error(t, err)

		var e cli.InvalidInputFlagError
		require.True(t, errors.As(err, &e))
	})
	t.Run("should successfully validate some inputs", func(t *testing.T) {
		testCases := []*cli.Config{
			{
				KeyLength: 24,
				EncText:   "someEncryptedText",
				PlainText: "somePlainText",
			},
			{
				KeyLength: 28,
				EncText:   "someEncryptedText",
				PlainText: "somePlainText",
			},
			{
				KeyLength: 32,
				EncText:   "someEncryptedText",
				PlainText: "somePlainText",
			},
		}
		for idx, testCase := range testCases {
			testCase := testCase
			t.Run(fmt.Sprintf("testcase n %d", idx), func(t *testing.T) {
				t.Parallel()
				require.NoError(t, testCase.Validate())
			})
		}
	})
}
