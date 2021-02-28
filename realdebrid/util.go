package realdebrid

import "errors"

// SelectLargestFile returns the file ID of the largest file in the torrent.
func SelectLargestFile(info TorrentInfo) (File, error) {
	var largestFile File
	largestSize := 0

	for _, file := range info.Files {
		if file.Bytes > largestSize {
			largestFile = file
			largestSize = file.Bytes
		}
	}
	if largestSize == 0 {
		return File{}, errors.New("couldn't find largest file in torrent info")
	}

	return largestFile, nil
}
