package premiumize_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/deflix-tv/go-debrid/premiumize"
)

// Night of the Living Dead, 1968, public domain (so legal to download, stream and share), from YTS
var (
	nightOfTheLivingDeadHash   = "50B7DAFB7137CBECF045F78E8EFBE4AC1A90D139"
	nightOfTheLivingDeadMagnet = "magnet:?xt=urn:btih:50B7DAFB7137CBECF045F78E8EFBE4AC1A90D139&dn=Night+of+the+Living+Dead+%281968%29+%5B720p%5D+%5BYTS.MX%5D&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2F9.rarbg.to%3A2710%2Fannounce&tr=udp%3A%2F%2Fp4p.arenabg.ch%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.cyberia.is%3A6969%2Fannounce&tr=http%3A%2F%2Fp4p.arenabg.com%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337%2Fannounce"
)

func TestClient(t *testing.T) {
	apiKey, ok := os.LookupEnv("PM_APIKEY")
	require.True(t, ok, "API key is missing from the environment")

	// Create client
	auth := premiumize.Auth{
		KeyOrToken: apiKey,
	}
	client := premiumize.NewClient(premiumize.DefaultClientOpts, auth, nil)

	ctx := context.Background()

	// Get account info
	accInfo, err := client.GetAccountInfo(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, accInfo.CustomerID)
	require.NotEmpty(t, accInfo.PremiumUntil)

	// Check cache
	cachedFiles, err := client.CheckCache(ctx, nightOfTheLivingDeadHash)
	require.NoError(t, err)
	fmt.Printf("Cached files: %+v\n", cachedFiles)
	// We assume that the torrent is cached when this test runs.
	// If it's not, download it to PM so that it's cached afterwards, or use a different torrent, then run the test again.
	require.NotEmpty(t, len(cachedFiles))
	cachedFile, found := cachedFiles[nightOfTheLivingDeadHash]
	require.True(t, found)
	require.NotEmpty(t, cachedFile.Filesize)

	// Create transfer
	createdTransfer, err := client.CreateTransfer(ctx, nightOfTheLivingDeadMagnet)
	require.NoError(t, err)
	fmt.Printf("ID: %v\n", createdTransfer.ID)
	require.NotEmpty(t, createdTransfer.ID)

	// List transfers
	transfers, err := client.ListTransfers(ctx)
	require.NoError(t, err)
	fmt.Printf("Transfers: %+v\n", transfers)
	require.NotEmpty(t, transfers)
	var transfer premiumize.Transfer
	for _, elem := range transfers {
		if elem.ID == createdTransfer.ID {
			transfer = elem
			break
		}
	}
	require.Equal(t, "finished", transfer.Status)

	// Create direct download link
	downloads, err := client.CreateDDL(ctx, transfer.Src)
	require.NoError(t, err)
	fmt.Printf("Downloads: %+v\n", downloads)
	require.NotEmpty(t, downloads)

	dl, err := premiumize.SelectLargestFile(downloads)
	require.NoError(t, err)
	fmt.Printf("Largest download: %+v\n", dl)
	require.NotEmpty(t, dl.Link)
}
