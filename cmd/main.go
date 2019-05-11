package main

import (
	"fmt"
	"log"

	"github.com/andream16/mitmcracker"
	"github.com/andream16/mitmcracker/internal/cracker"
	"github.com/andream16/mitmcracker/internal/repository/redis"
)

func main() {

	in, err := mitmcracker.New()
	if err != nil {
		log.Fatal(err)
	}

	repo, err := redis.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	cracker := cracker.New(
		in.PlainText,
		in.EncText,
		in.KeyLength,
		repo,
	)

	keys, err := cracker.Crack()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("found encoding key: %s & decoding key: %s", keys.EncKey, keys.DecKey))

}
