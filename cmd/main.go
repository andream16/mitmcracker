package main

import (
	"flag"
	"log"

	"github.com/andream16/mitmcracker/internal/cli"
)

func main() {

	var conf cli.Config

	flag.UintVar(&conf.KeyLength, "key", 0, "key length")
	flag.StringVar(&conf.PlainText, "plain", "", "known plain text")
	flag.StringVar(&conf.EncText, "encoded", "", "known encoded text")
	flag.StringVar(&conf.Storage.Type, "storage.type", "", "storage to be used")
	flag.StringVar(&conf.Storage.Address, "storage.address", "", "storage address")
	flag.StringVar(&conf.Storage.Password, "storage.password", "", "storage password")
	flag.IntVar(&conf.Storage.DB, "storage.db", 0, "storage db")

	if err := conf.Validate(); err != nil {
		log.Fatal(err)
	}

	//repo, err := newRepository(in)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//keys, err := cracker.New(
	//	in.PlainText,
	//	in.EncText,
	//	in.KeyLength,
	//	repo,
	//).Crack()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Println(fmt.Sprintf("found encoding key: %s & decoding key: %s", keys.Encode, keys.Decode))

}

//func newRepository(in *cli.Input) (repository.Repositorer, error) {
//
//	if in.Storage.Type == "disk" {
//		return redis.New(in.Storage.Address, in.Storage.Password, in.Storage.DB)
//	}
//
//	return memory.New(cracker.GetKeyNumber(in.KeyLength)), nil
//}
