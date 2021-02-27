package realdebrid

import "errors"

// SelectLargestFile returns the file ID of the largest file in the torrent.
func SelectLargestFile(info TorrentInfo) (int, error) {
	largestID := -1
	largestSize := 0

	for i, file := range info.Files {
		if file.Bytes > largestSize {
			largestID = i
			largestSize = file.Bytes
		}
	}
	if largestID == -1 {
		return 0, errors.New("couldn't find largest file in torrent info")
	}

	return largestID, nil
}
