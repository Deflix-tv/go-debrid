package alldebrid

import "errors"

// SelectLargestFile returns the link of the largest file in the torrent.
func SelectLargestFile(status Status) (Link, error) {
	var largestLink Link
	largestSize := 0

	for _, link := range status.Links {
		if link.Size > largestSize {
			largestLink = link
			largestSize = link.Size
		}
	}
	if largestSize == 0 {
		return Link{}, errors.New("couldn't find largest file in status")
	}

	return largestLink, nil
}
