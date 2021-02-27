package realdebrid_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/deflix-tv/go-debrid/realdebrid"
)

// Night of the Living Dead, 1968, public domain (so legal to download, stream and share), from YTS
var (
	nightOfTheLivingDeadHash   = "50B7DAFB7137CBECF045F78E8EFBE4AC1A90D139"
	nightOfTheLivingDeadMagnet = "magnet:?xt=urn:btih:50B7DAFB7137CBECF045F78E8EFBE4AC1A90D139&dn=Night+of+the+Living+Dead+%281968%29+%5B720p%5D+%5BYTS.MX%5D&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2F9.rarbg.to%3A2710%2Fannounce&tr=udp%3A%2F%2Fp4p.arenabg.ch%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.cyberia.is%3A6969%2Fannounce&tr=http%3A%2F%2Fp4p.arenabg.com%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337%2Fannounce"
)

func TestClient(t *testing.T) {
	apiKey, ok := os.LookupEnv("RD_APITOKEN")
	require.True(t, ok, "API token is missing from the environment")

	// Create client
	auth := realdebrid.Auth{
		KeyOrToken: apiKey,
	}
	client := realdebrid.NewClient(realdebrid.DefaultClientOpts, auth, nil)

	ctx := context.Background()

	// Get user
	user, err := client.GetUser(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, user.ID)
	require.NotEmpty(t, user.Premium)

	// Get instant availability
	availabilities, err := client.GetInstantAvailability(ctx, nightOfTheLivingDeadHash)
	require.NoError(t, err)
	fmt.Printf("Availabilities: %+v\n", availabilities)
	// We assume that the torrent is instantly available when this test runs.
	// If it's not, download it to RD so that it's available afterwards, or use a different torrent, then run the test again.
	require.NotEmpty(t, len(availabilities))
	availability, found := availabilities[nightOfTheLivingDeadHash]
	require.True(t, found)
	require.NotEmpty(t, len(availability))
	// File ID 1 is the main movie file
	availableFile, found := availability[1]
	require.True(t, found)
	require.NotEmpty(t, availableFile.Filesize)

	// Add magnet
	torrentID, err := client.AddMagnet(ctx, nightOfTheLivingDeadMagnet)
	require.NoError(t, err)
	fmt.Printf("ID: %v\n", torrentID)
	require.NotEmpty(t, torrentID)

	// Get torrent info
	info, err := client.GetTorrentInfo(ctx, torrentID)
	require.NoError(t, err)
	fmt.Printf("Torrent info: %+v\n", info)
	require.NotEmpty(t, info.ID)
	// Although one or more files of the torrent are instantly available, the torrent info is user-specific.
	// It's not regarded as downloaded if no file has been selected for download yet.
	require.Equal(t, "waiting_files_selection", info.Status)

	fileID, err := realdebrid.SelectLargestFile(info)
	require.NoError(t, err)

	// "Download" file
	err = client.SelectFiles(ctx, torrentID, fileID)
	require.NoError(t, err)

	// Get torrent info again.
	info, err = client.GetTorrentInfo(ctx, torrentID)
	require.NoError(t, err)
	fmt.Printf("Torrent info: %+v\n", info)
	require.NotEmpty(t, info.ID)
	// This time the torrent is seen as downloaded
	require.Equal(t, "downloaded", info.Status)
	require.NotEmpty(t, info.Links)

	// Get HTTP download URL
	dl, err := client.Unrestrict(ctx, info.Links[0], false)
	require.NoError(t, err)
	fmt.Printf("Download: %+v\n", dl)
	require.NotEmpty(t, dl.Download)
}
