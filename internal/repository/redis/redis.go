package redis

import (
	"log"

	"github.com/go-redis/redis"
)

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
	return r.Client.HMSet("dec", map[string]interface{}{
		key: cipherText,
	}).Err()
}

// InsertEnc adds a new entry in dec map.
func (r *Redis) InsertEnc(key, cipherText string) error {
	return r.Client.HMSet("enc", map[string]interface{}{
		key: cipherText,
	}).Err()
}

// FindKey returns the common key.
func (r *Redis) FindKey() (string, error) {
	log.Println("ok1")
	scanCmd := r.Client.HScan("dec", 0, "", 0)
	if scanCmd.Err() != nil {
		return "", scanCmd.Err()
	}
	log.Println("ok2")
	for {
		v := scanCmd.Iterator().Val()
		log.Println("ok3", v)
		getCmd := r.Client.HGet("enc", v)
		if getCmd.Err() != nil {
			return "", getCmd.Err()
		}
		log.Println("ok4")
		k, err := getCmd.Result()
		if err != nil {
			return "", getCmd.Err()
		}
		log.Println("ok5")
		if k != "" {
			log.Println("ok-yo")
			return k, nil
		}
		if !scanCmd.Iterator().Next() {
			log.Println("ok-exit")
			break
		}
		log.Println("ok-continue")
	}
	log.Println("ok-def")
	return "", nil
}
