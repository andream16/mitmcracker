package main

import (
	"fmt"
	"log"

	"github.com/andream16/mitmcracker"
	"github.com/andream16/mitmcracker/internal/cracker"
	"github.com/andream16/mitmcracker/internal/repository/memory"
)

func main() {

	in, err := mitmcracker.New()
	if err != nil {
		log.Fatal(err)
	}

	repo := memory.New(cracker.GetKeyNumber(in.KeyLength))

	cracker := cracker.New(
		in.PlainText,
		in.EncText,
		in.KeyLength,
		repo,
	)

	key, err := cracker.Crack()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("found key %s", key))

}
