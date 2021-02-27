package premiumize

// CreatedTransfer represents a transfer that has just been added to Premiumize.
type CreatedTransfer struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Download represents a direct download. If a transfer was created by adding a torrent, a Download is a file in that torrent.
type Download struct {
	Path            string `json:"path,omitempty"`
	Size            string `json:"size,omitempty"`
	Link            string `json:"link,omitempty"`
	StreamLink      string `json:"stream_link,omitempty"`
	TranscodeStatus string `json:"transcode_status,omitempty"`
}

// Transfer represents a transfer, like a torrent that has been added to Premiumize for a specific user.
type Transfer struct {
	ID string `json:"id,omitempty"`
	// Name of the torrent if the transfer was created by adding a torrent
	Name    string `json:"name,omitempty"`
	Message string `json:"message,omitempty"`
	// "waiting", "finished" etc.
	Status string `json:"status,omitempty"`
	// Download progress. Can be 0 for cached files that don't have to be downloaded.
	Progress float64 `json:"progress,omitempty"`
	// When the transfer was created by adding a torrent via magnet URL, then this is the magnet URL
	Src      string `json:"src,omitempty"`
	FolderID string `json:"folder_id,omitempty"`
	FileID   string `json:"file_id,omitempty"`
}

// AccountInfo contains info about a user account.
type AccountInfo struct {
	CustomerID   string  `json:"customer_id,omitempty"`
	PremiumUntil int     `json:"premium_until,omitempty"`
	LimitUsed    float64 `json:"limit_used,omitempty"`
	SpaceUsed    float64 `json:"space_used,omitempty"`
}

// CachedFile represents a file that's available in Premiumize's cache.
type CachedFile struct {
	Transcoded bool
	Filename   string
	Filesize   string
}
