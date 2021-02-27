# go-debrid

Go library for turning torrents into cached HTTP streams via debrid services like RealDebrid, AllDebrid and Premiumize

## Usage

```go
package main

import (
    "context"
    "fmt"

    "github.com/deflix-tv/go-debrid/realdebrid"
)

func main() {
    // Create new client
    auth := realdebrid.Auth{KeyOrToken: "123"}
    rd := realdebrid.NewClient(realdebrid.DefaultClientOpts, auth, nil)

    // We're using some info hashes of "Night of the Living Dead" from 1968, which is in the public domain
    infoHashes := []string{
        "50B7DAFB7137CBECF045F78E8EFBE4AC1A90D139",
        "11EA02584FA6351956F35671962AB46354D99060",
    }
    availabilities, _ := rd.GetInstantAvailability(context.Background(), infoHashes...)

    // Iterate through the available torrents and print their details
    for hash, availability := range availabilities {
        fmt.Printf("Hash: %v\nAvailability: %+v\n", hash, availability)
    }
}

```

For more detailed examples see [examples](examples).
