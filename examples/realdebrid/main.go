package main

import (
	"context"
	"fmt"

	"github.com/deflix-tv/go-debrid/realdebrid"
)

func main() {
	// Prepare client parameters
	opts := realdebrid.DefaultClientOpts
	auth := realdebrid.Auth{KeyOrToken: "123"}

	// Create new client
	rd := realdebrid.NewClient(opts, auth, nil)

	// We're using some info hashes of "Night of the Living Dead" from 1968, which is in the public domain
	infoHashes := []string{
		"50B7DAFB7137CBECF045F78E8EFBE4AC1A90D139",
		"11EA02584FA6351956F35671962AB46354D99060",
	}
	availabilities, err := rd.GetInstantAvailability(context.Background(), infoHashes...)
	if err != nil {
		panic(err)
	}
	if len(availabilities) == 0 {
		fmt.Println("None of the torrents are available")
	}

	// Iterate through the available torrents and print their details
	for hash, availability := range availabilities {
		fmt.Printf("Hash: %v\nAvailability: %+v\n", hash, availability)
	}
}
