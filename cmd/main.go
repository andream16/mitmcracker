package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/andream16/mitmcracker/internal/repository"

	"github.com/andream16/mitmcracker/internal/cli"
	"github.com/andream16/mitmcracker/internal/cracker"
	"github.com/andream16/mitmcracker/internal/decrypter"
	"github.com/andream16/mitmcracker/internal/encrypter"
	"github.com/andream16/mitmcracker/internal/formatter"
	"github.com/andream16/mitmcracker/internal/keycalculator"
	"github.com/andream16/mitmcracker/internal/perf"
	"github.com/andream16/mitmcracker/internal/repository/memory"
)

var version = "dev" // Set by build flags

func main() {
	// Setup structured logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	var (
		conf        cli.Config
		showVersion bool
	)

	flag.UintVar(&conf.KeyLength, "key", 0, "key length (24, 28, or 32 bits)")
	flag.StringVar(&conf.PlainText, "plain", "", "known plain text (16 hex digits)")
	flag.StringVar(&conf.EncText, "encoded", "", "known encoded text (16 hex digits)")
	flag.BoolVar(&showVersion, "version", false, "show version information")

	flag.Parse()

	if showVersion {
		fmt.Printf("mitmcracker version %s\n", version)
		return
	}

	slog.Info("Starting MITM cracker",
		"key_length", conf.KeyLength,
		"plain_text", conf.PlainText,
		"encoded_text", conf.EncText,
	)

	if err := conf.Validate(); err != nil {
		slog.Error("Configuration validation failed", "error", err)
		os.Exit(1)
	}

	slog.Info("Initializing cracker")
	cr, err := cracker.New(
		conf.KeyLength,
		conf.EncText,
		conf.PlainText,
		&memory.InMemo{
			KeyPairs: map[string]*repository.KeyPair{},
		},
		encrypter.DefaultEncrypt,
		decrypter.DefaultDecrypt,
		formatter.FastFormatter,
		keycalculator.DefaultCalculate,
		perf.DefaultMaxGoRoutineNumber,
		100000,
	)
	if err != nil {
		slog.Error("Failed to initialize cracker", "error", err)
		os.Exit(1)
	}

	slog.Info("Starting key cracking process")
	keyPair, found, err := cr.Crack(context.Background())
	if err != nil {
		slog.Error("Fatal error during key cracking", "error", err)
		os.Exit(1)
	}

	if !found {
		slog.Warn("No matching keys found")
		os.Exit(1)
	}

	slog.Info("Successfully found matching keys",
		"encode_key", keyPair.EncodeKey,
		"decode_key", keyPair.DecodeKey,
	)
	fmt.Printf("Found encoding key: %s & decoding key: %s\n", keyPair.EncodeKey, keyPair.DecodeKey)
}
