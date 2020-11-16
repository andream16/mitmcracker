package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/andream16/mitmcracker/internal/cli"
	"github.com/andream16/mitmcracker/internal/cracker"
	"github.com/andream16/mitmcracker/internal/decrypter"
	"github.com/andream16/mitmcracker/internal/encrypter"
	"github.com/andream16/mitmcracker/internal/formatter"
	"github.com/andream16/mitmcracker/internal/keycalculator"
	"github.com/andream16/mitmcracker/internal/perf"
	"github.com/andream16/mitmcracker/internal/repository/memory"
)

func main() {

	var conf cli.Config

	flag.UintVar(&conf.KeyLength, "key", 0, "key length")
	flag.StringVar(&conf.PlainText, "plain", "", "known plain text")
	flag.StringVar(&conf.EncText, "encoded", "", "known encoded text")

	flag.Parse()

	if err := conf.Validate(); err != nil {
		log.Fatal(err)
	}

	cr, err := cracker.New(
		conf.KeyLength,
		conf.EncText,
		conf.PlainText,
		&memory.InMemo{},
		encrypter.DefaultEncrypt,
		decrypter.DefaultDecrypt,
		formatter.FastFormatter,
		keycalculator.DefaultCalculate,
		perf.DefaultMaxGoRoutineNumber,
	)
	if err != nil {
		log.Fatalf("could not initialise cracker: %v", err)
	}

	keyPair, found, err := cr.Crack(context.Background())
	if err != nil {
		log.Fatalf("fatal error while finding the matching keys: %v", err)
	}

	if !found {
		log.Fatal("matching keys not found")
	}

	log.Println(fmt.Sprintf("found encoding key: %s & decoding key: %s", keyPair.EncodeKey, keyPair.DecodeKey))
}
