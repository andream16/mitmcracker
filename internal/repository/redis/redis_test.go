//+build integration

package redis

import (
	"testing"

	"github.com/go-redis/redis"
)

func TestRedis_InsertDec(t *testing.T) {
	r := newRedis(t)
	defer func() {
		err := r.Close()
		if err != redis.Nil {
			t.Fatalf("unexpected error %s", err)
		}
	}()
	err := r.InsertDec("key1", "val1")
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	cmd := r.Client.HGet("dec", "key1")
	if cmd.Err() != nil {
		t.Fatalf("unexpected error %s", cmd.Err())
	}
	v, err := cmd.Result()
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	if "val1" != v {
		t.Fatalf("expected %s, got %s", "val1", v)
	}
}

func TestRedis_InsertEnc(t *testing.T) {
	r := newRedis(t)
	defer func() {
		err := r.Close()
		if err != redis.Nil {
			t.Fatalf("unexpected error %s", err)
		}
	}()
	err := r.InsertEnc("key1", "val1")
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	cmd := r.Client.HGet("dec", "key1")
	if cmd.Err() != nil {
		t.Fatalf("unexpected error %s", cmd.Err())
	}
	v, err := cmd.Result()
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	if "val1" != v {
		t.Fatalf("expected %s, got %s", "val1", v)
	}
}

func TestRedis_FindKey(t *testing.T) {
	r := newRedis(t)
	defer func() {
		err := r.Close()
		if err != redis.Nil {
			t.Fatalf("unexpected error %s", err)
		}
	}()
	err := r.InsertEnc("cipher", "key1")
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	err = r.InsertDec("cipher", "key1")
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	key, err := r.FindKey()
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	if "ke1" != key {
		t.Fatalf("expected %s, got %s", "key1", key)
	}
}

func newRedis(t *testing.T) *Redis {
	cfg := &Config{
		Address: "localhost:6379",
	}
	r, err := New(cfg)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	return r
}
