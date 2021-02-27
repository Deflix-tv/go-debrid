# go-debrid

[![Go Reference](https://pkg.go.dev/badge/github.com/deflix-tv/go-debrid.svg)](https://pkg.go.dev/github.com/deflix-tv/go-debrid)

Go library for the public APIs of debrid services like [RealDebrid](https://real-debrid.com/), [AllDebrid](https://alldebrid.com/) and [Premiumize](https://www.premiumize.me/)

## Features

- Get user account info
- Get instant availability (cache) info for a link / torrent
- Add a link / torrent to a debrid service's downloads
- Get status info about a link / torrent that the debrid service is downloading / has downloaded
- Get the direct download link for a link / torrent after the debrid service has downloaded it

## Usage

The library consists of a root-level package which contains a cache interface and example cache implementation, as well as subpackages for the specific debrid services. Each service-specific subpackage contains both a legacy client (the client from `v0.1.0`), and a low level client whose methods match the public API endpoints. In the future a common client will be added that has a generic interface and is backed by service-specific clients.

Godoc:

- [RealDebrid](https://pkg.go.dev/github.com/deflix-tv/go-debrid/realdebrid)
- [AllDebrid](https://pkg.go.dev/github.com/deflix-tv/go-debrid/alldebrid)
- [Premiumize](https://pkg.go.dev/github.com/deflix-tv/go-debrid/premiumize)

### Example

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
