package redis

import (
	"errors"

	"github.com/go-redis/redis"
)

const (
	decKey = "dec"
	encKey = "enc"
)

var errNotFound = errors.New("key_not_found")

// Redis is a wrapper to redis.
type Redis struct {
	Client *redis.Client
}

// Config is the redis conf.
type Config struct {
	Address  string
	DB       int
	Password string
}

// Keys is the final result of FindKeys.
type Keys struct {
	EncKey string
	DecKey string
}

// New returns a new redis client.
func New(cfg *Config) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return &Redis{
		Client: client,
	}, client.Ping().Err()
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

// FindKeys returns the common key.
func (r *Redis) FindKeys() (*Keys, error) {
	iter := r.Client.HScan(decKey, 0, "", 0).Iterator()
	for iter.Next() {
		encK, err := r.get(encKey, iter.Val())
		if err != nil {
			if err == errNotFound {
				continue
			}
			return nil, err
		}
		if encK != "" {
			decK, err := r.get(decKey, iter.Val())
			if err != nil {
				if err == errNotFound {
					continue
				}
				return nil, err
			}
			return &Keys{
				EncKey: encK,
				DecKey: decK,
			}, nil
		}
	}
	return nil, iter.Err()
}

func (r *Redis) get(hmName, key string) (string, error) {
	cmd := r.Client.HGet(hmName, key)
	err := cmd.Err()
	if err != nil {
		if err == redis.Nil {
			return "", errNotFound
		}
		return "", err
	}
	res, err := cmd.Result()
	if err != nil {
		return "", err
	}
	return res, nil
}
