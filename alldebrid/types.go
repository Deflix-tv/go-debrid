package alldebrid

import (
	"errors"
)

var (
	// Bad request.
	// Corresponds to AllDebrid 400 status code.
	ErrorBadRequest = errors.New("bad request")
	// Authentication error.
	// Corresponds to AllDebrid 401 status code.
	ErrorUnauthorized = errors.New("unauthorized")
	// Too many requests hit the API too quickly, see https://docs.alldebrid.com/#rate-limiting.
	// Corresponds to AllDebrid 429 status code.
	ErrorTooManyRequests = errors.New("too many requests")
	// Something went wrong on Alldebrid's end.
	// Corresponds to AllDebrid 500, 502, 503 and 504 status codes.
	ErrorServerError = errors.New("server error")
)

// TODO: Add error vars for the endpoint-specific errors, where the error codes are part of the HTTP response body, like "LINK_HOST_NOT_SUPPORTED" when trying to unlock a link.

var errMap = map[int]error{
	400: ErrorBadRequest,
	401: ErrorUnauthorized,
	429: ErrorTooManyRequests,
	500: ErrorServerError,
	502: ErrorServerError,
	503: ErrorServerError,
	504: ErrorServerError,
}

// User represents an AllDebrid user.
type User struct {
	// User username
	Username string `json:"username"`
	// User email
	Email string `json:"email"`
	// true is premium, false if not
	IsPremium bool `json:"isPremium"`
	// true is user has active subscription, false if not
	IsSubscribed bool `json:"isSubscribed"`
	// true is account is in freedays trial, false if not
	IsTrial bool `json:"isTrial"`
	// 0 if user is not premium, or timestamp until user is premium
	PremiumUntil int `json:"premiumUntil"`
	// Language used by the user on Alldebrid, eg. 'en', 'fr'. Default to fr
	Lang string `json:"lang"`
	// Preferer TLD used by the user, eg. 'fr', 'es'. Default to fr
	PreferredDomain string `json:"preferedDomain"`
	// Number of fidelity points
	FidelityPoints int `json:"fidelityPoints"`
	// Remaining quotas for the limited hosts (in MB)
	LimitedHostersQuotas map[string]int `json:"limitedHostersQuotas"`
	// When in trial mode, remaining global traffic quota available (in MB)
	RemainingTrialQuota int `json:"remainingTrialQuota,omitempty"`
}

// Download represents an unlocked link.
type Download struct {
	// Requested link, simplified if it was not in canonical form
	Link string `json:"link,omitempty"`
	// Link's file filename
	Filename string `json:"filename,omitempty"`
	// Link host minified
	Host string `json:"host,omitempty"`
	// List of alternative links with other resolutions for some video links
	Streams []Stream `json:"streams,omitempty"`
	// Unused
	Paws bool `json:"paws,omitempty"`
	// Filesize of the link's file
	Filesize int `json:"filesize,omitempty"`
	// Generation ID
	ID string `json:"id,omitempty"`
	// Matched host main domain
	HostDomain string `json:"hostDomain,omitempty"`
	// Delayed ID if link need time to generate
	Delayed int `json:"delayed,omitempty"`
}

// Stream is an alternative stream with a different resolution than the original one.
type Stream struct {
	// Resolution, e.g. `480` if the resolution is 480p.
	Quality int `json:"quality,omitempty"`
	// E.g. "mp4"
	Ext string `json:"ext,omitempty"`
	// File size in bytes
	Filesize int    `json:"filesize,omitempty"`
	Name     string `json:"name,omitempty"`
	// Streamable direct link to the file
	Link string `json:"link,omitempty"`
	ID   string `json:"id,omitempty"`
}

// Magnet represents a magnet that was added to AllDebrid.
type Magnet struct {
	// Magnet sent
	Magnet string `json:"magnet,omitempty"`
	// Magnet filename, or 'noname' if could not parse it
	Name string `json:"name,omitempty"`
	// Magnet id, used to query status
	ID int `json:"id,omitempty"`
	// Magnet hash
	Hash string `json:"hash,omitempty"`
	// Magnet files size
	Size int `json:"size,omitempty"`
	// Whether the magnet is already available
	Ready bool `json:"ready,omitempty"`
}

// Status contains status info about a torrent that was previously uploaded to AllDebrid for a specific user.
type Status struct {
	// Magnet id
	ID int `json:"id,omitempty"`
	// Magnet filename
	Filename string `json:"filename,omitempty"`
	// Magnet filesize
	Size int `json:"size,omitempty"`
	// Status in plain English
	Status string `json:"status,omitempty"`
	// Status code
	StatusCode StatusCode `json:"statusCode,omitempty"`
	// Downloaded data so far
	Downloaded int `json:"downloaded,omitempty"`
	// Uploaded data so far
	Uploaded int `json:"uploaded,omitempty"`
	// Seeders count
	Seeders int `json:"seeders,omitempty"`
	// Download speed
	DownloadSpeed int `json:"downloadSpeed,omitempty"`
	// Upload speed
	UploadSpeed int `json:"uploadSpeed,omitempty"`
	// Timestamp of the date of the magnet upload
	UploadDate int `json:"uploadDate,omitempty"`
	// Timestamp of the date of the magnet completion
	CompletionDate int `json:"completionDate,omitempty"`
	// an array of link objects
	Links []Link `json:"links,omitempty"`
	// files array format
	Version int `json:"version,omitempty"`
}

// StatusCode indicates in which status an added torrent is.
type StatusCode int

const (
	StatusCode_InQueue StatusCode = iota
	StatusCode_Downloading
	StatusCode_CompressingMoving
	StatusCode_Uploading
	StatusCode_Ready
	StatusCode_UploadFail
	StatusCode_InternalErrorOnUnpacking
	StatusCode_NotDownloadedIn20Min
	StatusCode_FileTooBig
	StatusCode_InternalError
	StatusCode_DownloadTookMoreThan72h
	StatusCode_DeletedOnTheHosterWebsite
)

// Link represents a file in a torrent.
type Link struct {
	// Download link
	Link string `json:"link,omitempty"`
	// File name
	Filename string `json:"filename,omitempty"`
	// File size
	Size int `json:"size,omitempty"`
	// different format depending of version property
	Files []interface{} `json:"files,omitempty"`
}
