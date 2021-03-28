package realdebrid

import (
	"errors"
	"time"
)

var (
	// Bad request.
	// Corresponds to RealDebrid 400 status code.
	ErrorBadRequest = errors.New("bad request")
	// Expired, invalid.
	// Corresponds to RealDebrid 401 status code.
	ErrorBadToken = errors.New("bad token")
	// Account locked, not premium.
	// Corresponds to RealDebrid 403 status code.
	ErrorPermissionDenied = errors.New("permission denied")
	// Wrong parameter (invalid file id(s)) / Unknown ressource (invalid id)
	// Corresponds to RealDebrid 404 status code.
	ErrorInvalidID = errors.New("invalid ID")
	// Service unavailable.
	// Corresponds to RealDebrid 503 status code.
	ErrorServiceUnavailable = errors.New("service unavailable")
)

var errMap = map[int]error{
	400: ErrorBadRequest,
	401: ErrorBadToken,
	403: ErrorPermissionDenied,
	404: ErrorInvalidID,
	503: ErrorServiceUnavailable,
}

// User represents a RealDebrid user.
type User struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	// Fidelity points
	Points int `json:"points,omitempty"`
	// User language
	Locale string `json:"locale,omitempty"`
	Avatar string `json:"avatar,omitempty"`
	// "premium" or "free"
	Type string `json:"type,omitempty"`
	// seconds left as a Premium user
	Premium    int       `json:"premium,omitempty"`
	Expiration time.Time `json:"expiration,omitempty"`
}

// Download represents an unrestricted link.
type Download struct {
	ID       string `json:"id,omitempty"`
	Filename string `json:"filename,omitempty"`
	// Mime Type of the file, guessed by the file extension
	MimeType string `json:"mimeType,omitempty"`
	// Filesize in bytes, 0 if unknown
	Filesize int `json:"filesize,omitempty"`
	// Original link
	Link string `json:"link,omitempty"`
	// Host main domain
	Host string `json:"host,omitempty"`
	// Max Chunks allowed
	Chunks int `json:"chunks,omitempty"`
	// Disable / enable CRC check
	CRC int `json:"crc,omitempty"`
	// Generated link
	Download string `json:"download,omitempty"`
	// Is the file streamable on website
	Streamable int `json:"streamable,omitempty"`
}

// TorrentsInfo contains info about one element of a list of torrents that was added to RealDebrid for a specific user.
// It contains download info (progress, selected files) after one or more files of the torrent were selected to be downloaded.
// It's similar to TorrentInfo, but lacks some fields like OriginalFilename, OriginalBytes and Files.
type TorrentsInfo struct {
	ID       string `json:"id,omitempty"`
	Filename string `json:"filename,omitempty"`
	// SHA1 Hash of the torrent
	Hash string `json:"hash,omitempty"`
	// Size of selected files only
	Bytes int `json:"bytes,omitempty"`
	// Host main domain
	Host string `json:"host,omitempty"`
	// Split size of links
	Split int `json:"split,omitempty"`
	// Possible values: 0 to 100
	Progress int `json:"progress,omitempty"`
	// Current status of the torrent: magnet_error, magnet_conversion, waiting_files_selection, queued, downloading, downloaded, error, virus, compressing, uploading, dead
	Status string    `json:"status,omitempty"`
	Added  time.Time `json:"added,omitempty"`
	// Host URLs
	Links []string `json:"links,omitempty"`
	// !! Only present when finished, jsonDate
	Ended string `json:"ended,omitempty"`
	// !! Only present in "downloading", "compressing", "uploading" status
	Speed int `json:"speed,omitempty"`
	// !! Only present in "downloading", "magnet_conversion" status
	Seeders int `json:"seeders,omitempty"`
}

// TorrentInfo contains info about a specific torrent that was added to RealDebrid for a specific user.
// It contains download info (progress, selected files) after one or more files of the torrent were selected to be downloaded.
// It's similar to TorrentsInfo, but has some additional fields like OriginalFilename, OriginalBytes and Files.
type TorrentInfo struct {
	ID       string `json:"id,omitempty"`
	Filename string `json:"filename,omitempty"`
	// Original name of the torrent
	OriginalFilename string `json:"original_filename,omitempty"`
	// SHA1 Hash of the torrent
	Hash string `json:"hash,omitempty"`
	// Size of selected files only
	Bytes int `json:"bytes,omitempty"`
	// Total size of the torrent
	OriginalBytes int `json:"original_bytes,omitempty"`
	// Host main domain
	Host string `json:"host,omitempty"`
	// Split size of links
	Split int `json:"split,omitempty"`
	// Possible values: 0 to 100
	Progress int `json:"progress,omitempty"`
	// Current status of the torrent: magnet_error, magnet_conversion, waiting_files_selection, queued, downloading, downloaded, error, virus, compressing, uploading, dead
	Status string    `json:"status,omitempty"`
	Added  time.Time `json:"added,omitempty"`
	Files  []File    `json:"files,omitempty"`
	// Host URLs
	Links []string `json:"links,omitempty"`
	// !! Only present when finished, jsonDate
	Ended string `json:"ended,omitempty"`
	// !! Only present in "downloading", "compressing", "uploading" status
	Speed int `json:"speed,omitempty"`
	// !! Only present in "downloading", "magnet_conversion" status
	Seeders int `json:"seeders,omitempty"`
}

// File represents a file in a torrent.
type File struct {
	ID int `json:"id,omitempty"`
	// Path to the file inside the torrent, starting with "/"
	Path  string `json:"path,omitempty"`
	Bytes int    `json:"bytes,omitempty"`
	// 0 or 1
	Selected int `json:"selected,omitempty"`
}

// InstantAvailability maps torrent file IDs to their availability.
type InstantAvailability map[int]AvailableFile

// AvailableFile represents an instantly available file.
type AvailableFile struct {
	Filename string `json:"filename,omitempty"`
	Filesize int    `json:"filesize,omitempty"`
}
