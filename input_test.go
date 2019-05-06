package mitmcracker

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
)

func TestInput_validate(t *testing.T) {
	t.Run("should fail because key length is not supported", func(t *testing.T) {
		i := &Input{}
		err := i.validate()
		if errInvalidField != errors.Cause(err) {
			t.Fatalf("expected %s, got %v", errInvalidField, err)
		}
	})
	t.Run("should fail because key length is invalig", func(t *testing.T) {
		i := &Input{
			KeyLength: 10,
		}
		err := i.validate()
		if errInvalidField != errors.Cause(err) {
			t.Fatalf("expected %s, got %v", errInvalidField, err)
		}
	})
	t.Run("should fail because known encrypted text is empty", func(t *testing.T) {
		i := &Input{
			KeyLength: 24,
		}
		err := i.validate()
		if errInvalidField != errors.Cause(err) {
			t.Fatalf("expected %s, got %v", errInvalidField, err)
		}
	})
	t.Run("should fail because known plain text is empty", func(t *testing.T) {
		i := &Input{
			KeyLength: 28,
			EncText:   "someEncryptedText",
		}
		err := i.validate()
		if errInvalidField != errors.Cause(err) {
			t.Fatalf("expected %s, got %v", errInvalidField, err)
		}
	})
	t.Run("should fail because storage is not supported", func(t *testing.T) {
		i := &Input{
			KeyLength: 28,
			EncText:   "someEncryptedText",
			PlainText: "somePlainText",
			Storage:   "xyz",
		}
		err := i.validate()
		if errInvalidField != errors.Cause(err) {
			t.Fatalf("expected %s, got %v", errInvalidField, err)
		}
	})
	t.Run("should successfully validate some inputs", func(t *testing.T) {
		testCases := []*Input{
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
				err := testCase.validate()
				if err != nil {
					t.Fatalf("unexpected error %s", err)
				}
			})
		}
	})
}
