package alldebrid_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/deflix-tv/go-debrid/alldebrid"
)

// Night of the Living Dead, 1968, public domain (so legal to download, stream and share), from YTS
var (
	nightOfTheLivingDeadHash   = "50B7DAFB7137CBECF045F78E8EFBE4AC1A90D139"
	nightOfTheLivingDeadMagnet = "magnet:?xt=urn:btih:50B7DAFB7137CBECF045F78E8EFBE4AC1A90D139&dn=Night+of+the+Living+Dead+%281968%29+%5B720p%5D+%5BYTS.MX%5D&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2F9.rarbg.to%3A2710%2Fannounce&tr=udp%3A%2F%2Fp4p.arenabg.ch%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.cyberia.is%3A6969%2Fannounce&tr=http%3A%2F%2Fp4p.arenabg.com%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337%2Fannounce"
)

func TestClient(t *testing.T) {
	apiKey, ok := os.LookupEnv("AD_APIKEY")
	require.True(t, ok, "API key is missing from the environment")

	// Create client
	client := alldebrid.NewClient(alldebrid.DefaultClientOpts, apiKey, nil)

	ctx := context.Background()

	// Get user
	user, err := client.GetUser(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, user.Username)
	require.NotEmpty(t, user.PremiumUntil)

	// Get instant availability
	availabilities, err := client.GetInstantAvailability(ctx, nightOfTheLivingDeadHash)
	require.NoError(t, err)
	fmt.Printf("Availabilities: %+v\n", availabilities)
	// We assume that the torrent is instantly available when this test runs.
	// If it's not, download it to AD so that it's available afterwards, or use a different torrent, then run the test again.
	require.NotEmpty(t, len(availabilities))
	_, found := availabilities[nightOfTheLivingDeadHash]
	require.True(t, found)

	// Upload magnet
	magnet, err := client.UploadMagnet(ctx, nightOfTheLivingDeadMagnet)
	require.NoError(t, err)
	fmt.Printf("ID: %v\n", magnet.ID)
	require.NotEmpty(t, magnet.ID)

	// Get torrent info
	status, err := client.GetStatus(ctx, magnet.ID)
	require.NoError(t, err)
	fmt.Printf("Torrent status: %+v\n", status)
	require.NotEmpty(t, status.ID)
	require.Equal(t, alldebrid.StatusCode_Ready, status.StatusCode)

	link, err := alldebrid.SelectLargestFile(status)
	require.NoError(t, err)
	fmt.Printf("Largest link: %+v\n", link)
	require.NotEmpty(t, link.Link)

	// Get HTTP download URL
	dl, err := client.Unlock(ctx, link.Link)
	require.NoError(t, err)
	fmt.Printf("Download: %+v\n", dl)
	require.NotEmpty(t, dl.Link)
}
