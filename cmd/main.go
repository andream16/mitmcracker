package main

import (
	"fmt"
	"log"

	"github.com/andream16/mitmcracker"
)

func main() {

	in, err := mitmcracker.New()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(in.String())

}