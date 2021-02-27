package premiumize

import (
	"errors"
	"fmt"
	"strconv"
)

// SelectLargestFile returns the largest file in a slice of Download objects.
func SelectLargestFile(downloads []Download) (Download, error) {
	var result Download
	largestSize := 0

	for _, dl := range downloads {
		size, err := strconv.Atoi(dl.Size)
		if err != nil {
			return Download{}, fmt.Errorf("couldn't parse size: %w", err)
		}
		if size > largestSize {
			result = dl
			largestSize = size
		}
	}
	if largestSize == 0 {
		return Download{}, errors.New("couldn't find largest file in downloads")
	}

	return result, nil
}
