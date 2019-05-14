package redis

import (
	"github.com/andream16/mitmcracker/internal/repository"

	"github.com/go-redis/redis"
)

const (
	decKey = "dec"
	encKey = "enc"
)

var _ repository.Repositorer = (*Redis)(nil)

// Redis is a wrapper to redis.
type Redis struct {
	Client *redis.Client
}

// New returns a new redis client.
func New(
	address string,
	password string,
	db int,
) (*Redis, error) {
	opt := buildOptions(
		address,
		password,
		db,
	)
	client := redis.NewClient(opt)
	return &Redis{
		Client: client,
	}, client.Ping().Err()
}

func buildOptions(
	address string,
	password string,
	db int,
) *redis.Options {
	opt := &redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	}
	if opt.Addr == "" {
		opt.Addr = "localhost:6379"
	}
	return opt
}

// Close closes the connection to the client.
func (r *Redis) Close() error {
	return r.Client.ClientGetName().Err()
}

// InsertDec adds a new entry in dec map.
func (r *Redis) InsertDec(key, cipherText string) error {
	return r.Client.HMSet(decKey, map[string]interface{}{
		key: cipherText,
	}).Err()
}

// InsertEnc adds a new entry in dec map.
func (r *Redis) InsertEnc(key, cipherText string) error {
	return r.Client.HMSet(encKey, map[string]interface{}{
		key: cipherText,
	}).Err()
}

// FindKeys returns the keys used to encode and decode the message.
func (r *Redis) FindKeys() (*repository.Keys, error) {
	iter := r.Client.HScan(decKey, 0, "", 0).Iterator()
	for iter.Next() {
		encK, err := r.get(encKey, iter.Val())
		if err != nil {
			if err == repository.ErrNotFound {
				continue
			}
			return nil, err
		}
		if encK != "" {
			decK, err := r.get(decKey, iter.Val())
			if err != nil {
				if err == repository.ErrNotFound {
					continue
				}
				return nil, err
			}
			return &repository.Keys{
				Encode: encK,
				Decode: decK,
			}, nil
		}
	}
	if iter.Err() != nil {
		return nil, iter.Err()
	}
	return nil, repository.ErrNotFound
}

func (r *Redis) get(hmName, key string) (string, error) {
	cmd := r.Client.HGet(hmName, key)
	err := cmd.Err()
	if err != nil {
		if err == redis.Nil {
			return "", repository.ErrNotFound
		}
		return "", err
	}
	res, err := cmd.Result()
	if err != nil {
		return "", err
	}
	return res, nil
}
