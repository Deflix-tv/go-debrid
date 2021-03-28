package realdebrid

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

var zapDebridService = zap.String("debridService", "RealDebrid")

// ClientOptions are options for the client.
type ClientOptions struct {
	// Base URL for HTTP requests. This will also be used when making a request to a link that's read from a RealDebrid response by replacing its base URL.
	BaseURL string
	// Timeout for HTTP requests
	Timeout time.Duration
	// Extra headers to set for HTTP requests
	ExtraHeaders map[string]string
	// When setting this to true, the user's original IP address is read from Auth.IP and forwarded to RealDebrid for all POST requests.
	// Only required if the library is used in an app on a machine
	// whose outgoing IP is different from the machine that's going to request the cached file/stream URL.
	ForwardOriginIP bool
}

// DefaultClientOpts are ClientOptions with reasonable default values.
var DefaultClientOpts = ClientOptions{
	BaseURL: "https://api.real-debrid.com/rest/1.0",
	Timeout: 5 * time.Second,
}

// Auth carries authentication/authorization info for RealDebrid.
type Auth struct {
	// Long lasting API key or expiring OAuth2 access token
	KeyOrToken string
	// The user's original IP. Only required if ClientOptions.ForwardOriginIP is true.
	IP string
}

// Client represents a RealDebrid client.
type Client struct {
	opts       ClientOptions
	auth       Auth
	httpClient *http.Client
	logger     *zap.Logger
}

// NewClient returns a new RealDebrid client.
// The logger param can be nil.
func NewClient(opts ClientOptions, auth Auth, logger *zap.Logger) *Client {
	// Set default values
	if opts.BaseURL == "" {
		opts.BaseURL = DefaultClientOpts.BaseURL
	}
	if opts.Timeout == 0 {
		opts.Timeout = DefaultClientOpts.Timeout
	}
	if logger == nil {
		logger = zap.NewNop()
	}

	return &Client{
		opts: opts,
		auth: auth,
		httpClient: &http.Client{
			Timeout: opts.Timeout,
		},
		logger: logger,
	}
}

// GetUser fetches and returns the user object from RealDebrid.
func (c *Client) GetUser(ctx context.Context) (User, error) {
	c.logger.Debug("Getting user...", zapDebridService)

	resBytes, err := c.get(ctx, c.opts.BaseURL+"/user", nil)
	if err != nil {
		return User{}, fmt.Errorf("couldn't get user: %w", err)
	}
	user := User{}
	if err = json.Unmarshal(resBytes, &user); err != nil {
		return User{}, fmt.Errorf("couldn't unmarshal user: %w", err)
	}

	c.logger.Debug("Got user", zap.String("user", fmt.Sprintf("%+v", user)), zapDebridService)
	return user, nil
}

// Unrestrict unrestricts a hoster link.
// For torrents, the torrent must first be added to RealDebrid and a file selected for download, which then leads to such a hoster link.
// When remote is true, account sharing restrictions are lifted, but it requires separately purchased "sharing traffic".
func (c *Client) Unrestrict(ctx context.Context, link string, remote bool) (Download, error) {
	c.logger.Debug("Unrestricting link...", zapDebridService)

	data := url.Values{}
	data.Set("link", link)
	if remote {
		data.Set("remote", "1")
	}
	resBytes, err := c.post(ctx, c.opts.BaseURL+"/unrestrict/link", data)
	if err != nil {
		return Download{}, fmt.Errorf("couldn't unrestrict link: %w", err)
	}
	dl := Download{}
	if err = json.Unmarshal(resBytes, &dl); err != nil {
		return Download{}, fmt.Errorf("couldn't unmarshal download: %w", err)
	}

	c.logger.Debug("Unrestricted link", zap.String("download", fmt.Sprintf("%+v", dl)), zapDebridService)
	return dl, nil
}

// GetTorrentsInfo fetches and returns info about up to 100 torrents that were added to RealDebrid for a specific user.
// ActiveFirst leads to active torrents being the first in the returned list.
func (c *Client) GetTorrentsInfo(ctx context.Context, activeFirst bool) ([]TorrentsInfo, error) {
	c.logger.Debug("Getting torrents info...", zapDebridService)

	data := url.Values{}
	data.Set("offset", "0")
	data.Set("limit", "100")
	if activeFirst {
		data.Set("filter", "active")
	}
	resBytes, err := c.get(ctx, c.opts.BaseURL+"/torrents", data)
	if err != nil {
		return nil, fmt.Errorf("couldn't get torrents info: %w", err)
	}
	info := []TorrentsInfo{}
	if err = json.Unmarshal(resBytes, &info); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal torrents info: %w", err)
	}

	c.logger.Debug("Got torrents info", zap.String("torrentsInfo", fmt.Sprintf("%+v", info)), zapDebridService)
	return info, nil
}

// GetTorrentInfo fetches and returns info about a torrent that was added to RealDebrid for a specific user.
// The ID must be the one returned from RealDebrid when adding the torrent to RealDebrid.
func (c *Client) GetTorrentInfo(ctx context.Context, id string) (TorrentInfo, error) {
	c.logger.Debug("Getting torrent info...", zapDebridService)

	resBytes, err := c.get(ctx, c.opts.BaseURL+"/torrents/info/"+id, nil)
	if err != nil {
		return TorrentInfo{}, fmt.Errorf("couldn't get torrent info: %w", err)
	}
	info := TorrentInfo{}
	if err = json.Unmarshal(resBytes, &info); err != nil {
		return TorrentInfo{}, fmt.Errorf("couldn't unmarshal torrent info: %w", err)
	}

	c.logger.Debug("Got torrent info", zap.String("torrentInfo", fmt.Sprintf("%+v", info)), zapDebridService)
	return info, nil
}

// GetInstantAvailability fetches and returns info about the instant availability of a torrent.
func (c *Client) GetInstantAvailability(ctx context.Context, hashes ...string) (map[string]InstantAvailability, error) {
	c.logger.Debug("Getting instant availability...", zapDebridService)

	var hashParams string
	for _, hash := range hashes {
		hashParams += "/" + hash
	}
	resBytes, err := c.get(ctx, c.opts.BaseURL+"/torrents/instantAvailability"+hashParams, nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't get instant availability: %w", err)
	}
	availabilities := make(map[string]InstantAvailability, len(hashes))
	gjson.ParseBytes(resBytes).ForEach(func(key, value gjson.Result) bool {
		availableHash := key.String()
		// We check the original hash, so that our result has the same upper-/lowercase per hash
		for _, hash := range hashes {
			if strings.EqualFold(hash, availableHash) {
				availableHash = hash
				break
			}
		}
		availability := InstantAvailability{}
		value.Get("rd.0").ForEach(func(key, value gjson.Result) bool {
			availableFile := AvailableFile{}
			if err := json.Unmarshal([]byte(value.Raw), &availableFile); err != nil {
				c.logger.Error("Couldn't unmarshal available file", zap.Error(err), zap.String("availableFile", value.Raw), zapDebridService)
				return true
			}
			availability[int(key.Int())] = availableFile
			// Continue ForEach
			return true
		})
		if len(availability) > 0 {
			availabilities[availableHash] = availability
		}
		// Continue ForEach
		return true
	})

	c.logger.Debug("Got instant availability", zap.String("availabilities", fmt.Sprintf("%+v", availabilities)), zapDebridService)
	return availabilities, nil
}

// AddMagnet adds a torrent to RealDebrid via magnet URL.
func (c *Client) AddMagnet(ctx context.Context, magnet string) (string, error) {
	c.logger.Debug("Adding magnet...", zapDebridService)

	data := url.Values{}
	data.Set("magnet", magnet)
	resBytes, err := c.post(ctx, c.opts.BaseURL+"/torrents/addMagnet", data)
	if err != nil {
		return "", fmt.Errorf("couldn't add magnet: %w", err)
	}
	id := gjson.GetBytes(resBytes, "id").String()

	c.logger.Debug("Added magnet", zap.String("id", id), zapDebridService)
	return id, nil
}

// SelectFiles starts downloading the selected files from a torrent that was previously added to RealDebrid for the specific user.
func (c *Client) SelectFiles(ctx context.Context, torrentID string, fileIDs ...int) error {
	c.logger.Debug("Selecting files...", zapDebridService)

	data := url.Values{}
	idStrings := make([]string, len(fileIDs))
	for i, id := range fileIDs {
		idStrings[i] = strconv.Itoa(id)
	}
	idString := strings.Join(idStrings, ",")
	data.Set("files", idString)
	_, err := c.post(ctx, c.opts.BaseURL+"/torrents/selectFiles/"+torrentID, data)
	if err != nil {
		return fmt.Errorf("couldn't select files: %w", err)
	}

	c.logger.Debug("Selected files", zapDebridService)
	return nil
}
