package memory

import "testing"

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

		k := m.FindKey()
		if k == "" {
			t.Fatalf("expected %s, got an empty string", key)
		}

	})

	t.Run("should not find a common cipherText", func(t *testing.T) {

		cipherText := "AUSFBVA"
		key := "1234"

		m := &InMemo{}
		m.InsertEnc(key, cipherText)
		m.InsertDec(key, "AAAAA")

		k := m.FindKey()
		if k != "" {
			t.Fatalf("expected an empty string, got %s", k)
		}

	})

}
