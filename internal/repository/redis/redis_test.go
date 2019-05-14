
package redis

import (
	"testing"

	"github.com/andream16/mitmcracker/internal/repository"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

func TestRedis_InsertDec(t *testing.T) {
	r := newRedis(t)
	defer func() {
		cleanUp(r, t, decKey)
		err := r.Close()
		if err != redis.Nil {
			t.Fatalf("unexpected error %s", err)
		}
	}()
	err := r.InsertDec("cipherDec", "keyDec")
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	v, err := r.get(decKey, "cipherDec")
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	if "keyDec" != v {
		t.Fatalf("expected %s, got %s", "keyDec", v)
	}
}

func TestRedis_InsertEnc(t *testing.T) {
	r := newRedis(t)
	defer func() {
		cleanUp(r, t, encKey)
		err := r.Close()
		if err != redis.Nil {
			t.Fatalf("unexpected error %s", err)
		}
	}()
	err := r.InsertEnc("cipherEnc", "keyEnc")
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	v, err := r.get(encKey, "cipherEnc")
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	if "keyEnc" != v {
		t.Fatalf("expected %s, got %s", "keyEnc", v)
	}
}

func TestRedis_FindKey(t *testing.T) {
	t.Run("should find two keys for same cipher text", func(t *testing.T) {
		r := newRedis(t)
		defer func() {
			cleanUp(r, t, encKey, decKey)
			err := r.Close()
			if err != redis.Nil {
				t.Fatalf("unexpected error %s", err)
			}
		}()
		err := r.InsertEnc("cipher", "encKey")
		if err != nil {
			t.Fatalf("unexpected error %s", err)
		}
		err = r.InsertDec("cipher", "decKey")
		if err != nil {
			t.Fatalf("unexpected error %s", err)
		}
		keys, err := r.FindKeys()
		if err != nil {
			t.Fatalf("unexpected error %s", err)
		}
		if "encKey" != keys.Encode {
			t.Fatalf("expected %s, got %s", "encKey", keys.Encode)
		}
		if "decKey" != keys.Decode {
			t.Fatalf("expected %s, got %s", "decKey", keys.Decode)
		}
	})
	t.Run("should return not found error", func(t *testing.T) {
		r := newRedis(t)
		defer func() {
			cleanUp(r, t, encKey, decKey)
			err := r.Close()
			if err != redis.Nil {
				t.Fatalf("unexpected error %s", err)
			}
		}()
		err := r.InsertEnc("c", "e")
		if err != nil {
			t.Fatalf("unexpected error %s", err)
		}
		err = r.InsertDec("a", "d")
		if err != nil {
			t.Fatalf("unexpected error %s", err)
		}
		_, err = r.FindKeys()
		if repository.ErrNotFound != errors.Cause(err) {
			t.Fatalf("expected error not found, got %s", err)
		}
	})
}

func newRedis(t *testing.T) *Redis {
	r, err := New("", "", 0)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	return r
}

func cleanUp(r *Redis, t *testing.T, keys ...string) {
	err := r.Client.Del(keys...).Err()
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
}
