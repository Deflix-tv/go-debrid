# go-debrid

Go library for turning torrents into cached HTTP streams via debrid services like RealDebrid, AllDebrid and Premiumize

## Usage

```go
package main

import (
    "context"
    "fmt"

    "github.com/deflix-tv/go-debrid"
    "github.com/deflix-tv/go-debrid/realdebrid"
    "go.uber.org/zap"
)

func main() {
    // Create new client
    rd, err := realdebrid.NewClient(realdebrid.DefaultClientOpts, debrid.NewInMemoryCache(), debrid.NewInMemoryCache(), zap.NewNop())
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

    // Iterate through the available info hashes and print them
    for _, availableInfoHash := range availableInfoHashes {
        fmt.Printf("Available info_hash: %v\n", availableInfoHash)
    }
}
```

For more detailed examples see [examples](examples).
