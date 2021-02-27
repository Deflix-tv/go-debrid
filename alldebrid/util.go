package alldebrid

import "errors"

// SelectLargestFile returns the link of the largest file in the torrent.
func SelectLargestFile(status Status) (string, error) {
	var largestLink string
	largestSize := 0

	for _, link := range status.Links {
		if link.Size > largestSize {
			largestLink = link.Link
			largestSize = link.Size
		}
	}
	if largestLink == "" {
		return "", errors.New("couldn't find largest file in status")
	}

	return largestLink, nil
}
