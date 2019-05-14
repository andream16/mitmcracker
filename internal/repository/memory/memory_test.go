package memory

import (
	"testing"

	"github.com/andream16/mitmcracker/internal/repository"
	"github.com/pkg/errors"
)

func TestInMemo_InsertDec(t *testing.T) {

	t.Run("successfully add an entry", func(t *testing.T) {

		cipherText := "AUSFBVA"
		key := "1234"

		m := &InMemo{}
		m.InsertDec(key, cipherText)
		v, ok := m.Dec[cipherText]
		if !ok {
			t.Fatal("expected true, got false")
		}
		if key != v {
			t.Fatalf("expected %s, got %s", "someInput", v)
		}
	})

}

func TestInMemo_InsertEnc(t *testing.T) {

	t.Run("successfully add an entry", func(t *testing.T) {

		cipherText := "AUSFBVA"
		key := "1234"

		m := &InMemo{}
		m.InsertEnc(key, cipherText)
		v, ok := m.Enc[cipherText]
		if !ok {
			t.Fatal("expected true, got false")
		}
		if key != v {
			t.Fatalf("expected %s, got %s", "someInput", v)
		}
	})

}

func TestFindKey(t *testing.T) {

	t.Run("should find a common cipherText", func(t *testing.T) {

		cipherText := "AUSFBVA"
		key := "1234"

		m := &InMemo{}
		m.InsertEnc(key, cipherText)
		m.InsertDec(key, cipherText)

		keys, err := m.FindKeys()
		if err != nil {
			t.Fatalf("unexpected error %s", err)
		}
		if keys == nil {
			t.Fatal("result is nil")
		}
		if key != keys.Encode {
			t.Fatalf("expected key %s, got %s", key, keys.Encode)
		}
		if key != keys.Decode {
			t.Fatalf("expected key %s, got %s", key, keys.Decode)
		}
	})

	t.Run("should not find a common cipherText", func(t *testing.T) {

		cipherText := "AUSFBVA"
		key := "1234"

		m := &InMemo{}
		m.InsertEnc(key, cipherText)
		m.InsertDec(key, "AAAAA")

		k, err := m.FindKeys()
		if repository.ErrNotFound != errors.Cause(err) {
			t.Fatalf("expected %s, got %s", repository.ErrNotFound, err)
		}
		if k != nil {
			t.Fatalf("expected an empty string, got %s", k)
		}

	})

}
