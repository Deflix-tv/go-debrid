package debrid

import (
	"context"
	"time"
)

type Info struct {
	// TODO: Add common fields
}

type Download struct {
	// TODO: Add common fields
}

// Adapter is the interface for common debrid client functionality.
type Adapter interface {
	// GetUser fetches and returns the user of the debrid service.
	// The returned time is the date until when the user's premium subscription is valid.
	GetUser(context.Context) (time.Time, error)
	// GetInstantAvailability checks if files are cached on the debrid service.
	// The string parameters are torrent info hashes.
	// The returned map contains the info hashes of the torrents that are instantly available
	// and maps them to the instantly available files.
	// Some debrid services cache *all* files when caching a torrent, some don't.
	GetInstantAvailability(context.Context, ...string) (map[string]map[int]struct{}, error)
	// AddMagnet adds a torrent to the debrid service via magnet URL.
	// The debrid service either has it already in the cache, or starts downloading the torrent.
	// The returned string is the ID to be used for turning the torrent into a direct download later.
	// Calling it multiple times for the same magnet URL doesn't lead to an error.
	AddMagnet(context.Context, string) (string, error)
	// GetInfo fetches and returns information about a previously added torrent.
	GetInfo(context.Context, string) (Info, error)
	// CreateDDL creates a Download for the ID returned by AddMagnet.
	CreateDDL(context.Context, string) (map[int]Download, error)
}
