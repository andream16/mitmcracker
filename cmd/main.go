package main

import (
	"fmt"
	"log"

	"github.com/andream16/mitmcracker"
	"github.com/andream16/mitmcracker/internal/cracker"
	"github.com/andream16/mitmcracker/internal/repository"
	"github.com/andream16/mitmcracker/internal/repository/memory"
	"github.com/andream16/mitmcracker/internal/repository/redis"
)

func main() {

	in, err := mitmcracker.New()
	if err != nil {
		log.Fatal(err)
	}

	repo, err := newRepository(in)
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

	fmt.Println(fmt.Sprintf("found encoding key: %s & decoding key: %s", keys.Encode, keys.Decode))

}

func newRepository(in *mitmcracker.Input) (repository.Repositorer, error) {

	if in.Storage.Type == "disk" {
		return redis.New(in.Storage.Address, in.Storage.Password, in.Storage.DB)
	}

	return memory.New(cracker.GetKeyNumber(in.KeyLength)), nil

}
