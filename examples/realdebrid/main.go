package main

import (
	"context"
	"fmt"

	"github.com/deflix-tv/go-debrid"
	"github.com/deflix-tv/go-debrid/realdebrid"
	"go.uber.org/zap"
)

func main() {
	// Prepare client parameters
	opts := realdebrid.DefaultClientOpts
	tokenCache := debrid.NewInMemoryCache()
	availabilityCache := debrid.NewInMemoryCache()
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	// Create new client
	rd, err := realdebrid.NewClient(opts, tokenCache, availabilityCache, logger)
	if err != nil {
		panic(err)
	}

	auth := realdebrid.Auth{KeyOrToken: "123"}
	// We're using some info hashes of "Night of the Living Dead" from 1968, which is in the public domain
	infoHashes := []string{
		"50B7DAFB7137CBECF045F78E8EFBE4AC1A90D139",
		"11EA02584FA6351956F35671962AB46354D99060",
	}
	availableInfoHashes := rd.CheckInstantAvailability(context.Background(), auth, infoHashes...)
	if len(availableInfoHashes) == 0 {
		fmt.Println("None of the info hashes are available")
	}

	// Iterate through the available info hashes and print them
	for _, availableInfoHash := range availableInfoHashes {
		fmt.Printf("Available info_hash: %v\n", availableInfoHash)
	}
}
