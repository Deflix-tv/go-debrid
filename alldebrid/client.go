package alldebrid

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

var zapDebridService = zap.String("debridService", "AllDebrid")

// ClientOptions are options for the client.
type ClientOptions struct {
	// Base URL for HTTP requests. This will also be used when making a request to a link that's read from a AllDebrid response by replacing its base URL.
	BaseURL string
	// Timeout for HTTP requests
	Timeout time.Duration
	// Extra headers to set for HTTP requests
	ExtraHeaders map[string]string
}

// DefaultClientOpts are ClientOptions with reasonable default values.
var DefaultClientOpts = ClientOptions{
	BaseURL: "https://api.alldebrid.com/v4",
	Timeout: 5 * time.Second,
}

// Client represents a AllDebrid client.
type Client struct {
	opts       ClientOptions
	apiKey     string
	httpClient *http.Client
	logger     *zap.Logger
}

// NewClient returns a new AllDebrid client.
// The logger param can be nil.
func NewClient(opts ClientOptions, apiKey string, logger *zap.Logger) *Client {
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
		opts:   opts,
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: opts.Timeout,
		},
		logger: logger,
	}
}

// GetUser fetches and returns the user object from AllDebrid.
func (c *Client) GetUser(ctx context.Context) (User, error) {
	c.logger.Debug("Getting user...", zapDebridService)

	resBytes, err := c.get(ctx, c.opts.BaseURL+"/user")
	if err != nil {
		return User{}, fmt.Errorf("couldn't get user: %w", err)
	}
	if gjson.GetBytes(resBytes, "status").String() != "success" {
		errorCode := gjson.GetBytes(resBytes, "error.message").String()
		return User{}, fmt.Errorf("got error response from AllDebrid: %v", errorCode)
	}
	userJSON := gjson.GetBytes(resBytes, "data.user").Raw
	user := User{}
	if err = json.Unmarshal([]byte(userJSON), &user); err != nil {
		return User{}, fmt.Errorf("couldn't unmarshal user: %w", err)
	}

	c.logger.Debug("Got user", zap.String("user", fmt.Sprintf("%+v", user)), zapDebridService)
	return user, nil
}

// Unlock unlocks a link.
// For torrents, the torrent must first be added to AllDebrid, which then leads to such a hoster link (either instantly or after it was downloaded by AllDebrid).
func (c *Client) Unlock(ctx context.Context, link string) (Download, error) {
	c.logger.Debug("Unlocking link...", zapDebridService)

	resBytes, err := c.get(ctx, c.opts.BaseURL+"/link/unlock?link="+url.QueryEscape(link))
	if err != nil {
		return Download{}, fmt.Errorf("couldn't unlock link: %w", err)
	}
	if gjson.GetBytes(resBytes, "status").String() != "success" {
		errorCode := gjson.GetBytes(resBytes, "error.message").String()
		return Download{}, fmt.Errorf("got error response from AllDebrid: %v", errorCode)
	}
	downloadJSON := gjson.GetBytes(resBytes, "data").Raw
	dl := Download{}
	if err = json.Unmarshal([]byte(downloadJSON), &dl); err != nil {
		return Download{}, fmt.Errorf("couldn't unmarshal download: %w", err)
	}

	c.logger.Debug("Unlocked link", zap.String("download", fmt.Sprintf("%+v", dl)), zapDebridService)
	return dl, nil
}

// UploadMagnet adds a torrent to AllDebrid via magnet URL.
// The magnet string can actually also be a hash.
func (c *Client) UploadMagnet(ctx context.Context, magnet string) (Magnet, error) {
	c.logger.Debug("Uploading magnet...", zapDebridService)

	data := url.Values{}
	data.Set("magnets[]", magnet)
	resBytes, err := c.post(ctx, c.opts.BaseURL+"/magnet/upload", data)
	if err != nil {
		return Magnet{}, fmt.Errorf("couldn't upload magnet: %w", err)
	}
	if gjson.GetBytes(resBytes, "status").String() != "success" {
		errorCode := gjson.GetBytes(resBytes, "error.message").String()
		return Magnet{}, fmt.Errorf("got error response from AllDebrid: %v", errorCode)
	}
	magnetJSON := gjson.GetBytes(resBytes, "data.magnets.0").Raw
	m := Magnet{}
	if err = json.Unmarshal([]byte(magnetJSON), &m); err != nil {
		return Magnet{}, fmt.Errorf("couldn't unmarshal magnet: %w", err)
	}

	c.logger.Debug("Uploaded magnet", zap.String("magnet", fmt.Sprintf("%+v", m)), zapDebridService)
	return m, nil
}

// GetStatus fetches and returns the status of all torrents that were added to AllDebrid for a specific user.
// The ID must be the one returned from AllDebrid when adding the torrent to AllDebrid.
func (c *Client) GetStatus(ctx context.Context) ([]Status, error) {
	c.logger.Debug("Getting status...", zapDebridService)

	resBytes, err := c.get(ctx, c.opts.BaseURL+"/magnet/status")
	if err != nil {
		return nil, fmt.Errorf("couldn't get status: %w", err)
	}
	if gjson.GetBytes(resBytes, "status").String() != "success" {
		errorCode := gjson.GetBytes(resBytes, "error.message").String()
		return nil, fmt.Errorf("got error response from AllDebrid: %v", errorCode)
	}
	statusJSON := gjson.GetBytes(resBytes, "data.magnets").Raw
	status := []Status{}
	if err = json.Unmarshal([]byte(statusJSON), &status); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal status: %w", err)
	}

	c.logger.Debug("Got status", zap.String("status", fmt.Sprintf("%+v", status)), zapDebridService)
	return status, nil
}

// GetStatusByID fetches and returns the status of a specific torrent that was added to AllDebrid for a specific user.
// The ID must be the one returned from AllDebrid when adding the torrent to AllDebrid.
func (c *Client) GetStatusByID(ctx context.Context, id int) (Status, error) {
	c.logger.Debug("Getting status by ID...", zapDebridService)

	resBytes, err := c.get(ctx, c.opts.BaseURL+"/magnet/status?id="+strconv.Itoa(id))
	if err != nil {
		return Status{}, fmt.Errorf("couldn't get status by ID: %w", err)
	}
	if gjson.GetBytes(resBytes, "status").String() != "success" {
		errorCode := gjson.GetBytes(resBytes, "error.message").String()
		return Status{}, fmt.Errorf("got error response from AllDebrid: %v", errorCode)
	}
	statusJSON := gjson.GetBytes(resBytes, "data.magnets").Raw
	status := Status{}
	if err = json.Unmarshal([]byte(statusJSON), &status); err != nil {
		return Status{}, fmt.Errorf("couldn't unmarshal status by ID: %w", err)
	}

	c.logger.Debug("Got status by ID", zap.String("status", fmt.Sprintf("%+v", status)), zapDebridService)
	return status, nil
}

// DeleteMagnet deletes a magnet from the user's magnets.
// The ID must be the one returned from AllDebrid when adding the magnet or getting status info about it.
func (c *Client) DeleteMagnet(ctx context.Context, id int) error {
	c.logger.Debug("Deleting magnet...", zapDebridService)

	resBytes, err := c.get(ctx, c.opts.BaseURL+"/magnet/delete?id="+strconv.Itoa(id))
	if err != nil {
		return fmt.Errorf("couldn't delete magnet: %w", err)
	}
	if gjson.GetBytes(resBytes, "status").String() != "success" {
		errorCode := gjson.GetBytes(resBytes, "error.message").String()
		return fmt.Errorf("got error response from AllDebrid: %v", errorCode)
	}

	c.logger.Debug("Deleted magnet", zapDebridService)
	return nil
}

// GetInstantAvailability fetches and returns info about the instant availability of a torrent.
// The hashes can actually also be magnet URLs.
// The returned map contains the hashes / magnet URLs of the torrents that are instantly available.
func (c *Client) GetInstantAvailability(ctx context.Context, hashes ...string) (map[string]struct{}, error) {
	c.logger.Debug("Getting instant availability...", zapDebridService)

	data := url.Values{"magnets[]": hashes}
	resBytes, err := c.post(ctx, c.opts.BaseURL+"/magnet/instant", data)
	if err != nil {
		return nil, fmt.Errorf("couldn't get instant availability: %w", err)
	}
	if gjson.GetBytes(resBytes, "status").String() != "success" {
		errorCode := gjson.GetBytes(resBytes, "error.message").String()
		return nil, fmt.Errorf("got error response from AllDebrid: %v", errorCode)
	}
	availabilities := make(map[string]struct{}, len(hashes))
	gjson.GetBytes(resBytes, "data.magnets").ForEach(func(key, value gjson.Result) bool {
		if value.Get("instant").Bool() {
			availabilities[value.Get("magnet").String()] = struct{}{}
		}
		return true
	})

	c.logger.Debug("Got instant availability", zap.String("availabilities", fmt.Sprintf("%+v", availabilities)), zapDebridService)
	return availabilities, nil
}
