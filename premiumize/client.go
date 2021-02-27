package premiumize

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

var zapDebridService = zap.String("debridService", "Premiumize")

// ClientOptions are options for the client.
type ClientOptions struct {
	// Base URL for HTTP requests. This will also be used when making a request to a link that's read from a Premiumize response by replacing its base URL.
	BaseURL string
	// Timeout for HTTP requests
	Timeout time.Duration
	// Extra headers to set for HTTP requests
	ExtraHeaders map[string]string
	// When setting this to true, the user's original IP address is read from Auth.IP and forwarded to Premiumize when creating a direct download links.
	// Only required if the library is used in an app on a machine
	// whose outgoing IP is different from the machine that's going to request the cached file/stream URL.
	ForwardOriginIP bool
}

// DefaultClientOpts are ClientOptions with reasonable default values.
var DefaultClientOpts = ClientOptions{
	BaseURL: "https://www.premiumize.me/api",
	Timeout: 5 * time.Second,
}

// Auth carries authentication/authorization info for Premiumize.
type Auth struct {
	// Long lasting API key or expiring OAuth2 access token
	KeyOrToken string
	// Flag for indicating whether KeyOrToken is a key (false) or token (true).
	OAuth2 bool
	// The user's original IP. Only required if ClientOptions.ForwardOriginIP is true.
	IP string
}

// Client represents a Premiumize client.
type Client struct {
	opts       ClientOptions
	auth       Auth
	httpClient *http.Client
	logger     *zap.Logger
}

// NewClient returns a new Premiumize client.
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

// CreateTransfer creates a transfer.
// The source can be an HTTP(S) link to a supported container file, website or magnet link.
// Transfers that are created this way will appear in the transfer list.
func (c *Client) CreateTransfer(ctx context.Context, source string) (CreatedTransfer, error) {
	c.logger.Debug("Creating transfer...", zapDebridService)

	data := url.Values{}
	data.Set("src", source)
	resBytes, err := c.post(ctx, c.opts.BaseURL+"/transfer/create", data, true)
	if err != nil {
		return CreatedTransfer{}, fmt.Errorf("couldn't create transfer: %w", err)
	}
	if gjson.GetBytes(resBytes, "status").String() != "success" {
		message := gjson.GetBytes(resBytes, "message").String()
		return CreatedTransfer{}, fmt.Errorf("got error response from Premiumize: %v", message)
	}
	tf := CreatedTransfer{}
	if err = json.Unmarshal(resBytes, &tf); err != nil {
		return CreatedTransfer{}, fmt.Errorf("couldn't unmarshal added transfer: %w", err)
	}

	c.logger.Debug("Created transfer", zap.String("createdTransfer", fmt.Sprintf("%+v", tf)), zapDebridService)
	return tf, nil
}

// CreateDDL creates direct download links.
// The source can be an HTTP(S) link to a supported container file, website or magnet link.
// The creation will only work if the file is cached on Premiumize or if a transfer for the file has been created before and the transfer finished downloading (to Premiumize).
// If the source contains multiple files, each file is an element in the slice of Download objects.
func (c *Client) CreateDDL(ctx context.Context, source string) ([]Download, error) {
	c.logger.Debug("Creating direct download link...", zapDebridService)

	data := url.Values{}
	data.Set("src", source)
	// Premiumize asks for the original IP only for directdl requests
	if c.opts.ForwardOriginIP {
		if c.auth.IP == "" {
			return nil, errors.New("auth.IP is empty but client is configured to forward the user's original IP")
		}
		data.Add("download_ip", c.auth.IP)
	}
	resBytes, err := c.post(ctx, c.opts.BaseURL+"/transfer/directdl", data, true)
	if err != nil {
		return nil, fmt.Errorf("couldn't create direct download link: %w", err)
	}
	if gjson.GetBytes(resBytes, "status").String() != "success" {
		message := gjson.GetBytes(resBytes, "message").String()
		return nil, fmt.Errorf("got error response from Premiumize: %v", message)
	}
	downloadsJSON := gjson.GetBytes(resBytes, "content").Raw
	downloads := []Download{}
	if err = json.Unmarshal([]byte(downloadsJSON), &downloads); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal added transfer: %w", err)
	}

	c.logger.Debug("Created direct download link", zap.String("downloads", fmt.Sprintf("%+v", downloads)), zapDebridService)
	return downloads, nil
}

// ListTransfers fetches and returns all transfers that were previously added to Premiumize for a specific user.
// This doesn't include downloads that were created with CreateDDL without having been added via CreateTransfer.
func (c *Client) ListTransfers(ctx context.Context) ([]Transfer, error) {
	c.logger.Debug("Listing transfers...", zapDebridService)

	resBytes, err := c.get(ctx, c.opts.BaseURL+"/transfer/list", nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't list transfers: %w", err)
	}
	if gjson.GetBytes(resBytes, "status").String() != "success" {
		message := gjson.GetBytes(resBytes, "message").String()
		return nil, fmt.Errorf("got error response from Premiumize: %v", message)
	}
	transferJSON := gjson.GetBytes(resBytes, "transfers").Raw
	transfers := []Transfer{}
	if err = json.Unmarshal([]byte(transferJSON), &transfers); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal transfer list: %w", err)
	}

	c.logger.Debug("Created direct download link", zap.String("transfers", fmt.Sprintf("%+v", transfers)), zapDebridService)
	return transfers, nil
}

// GetAccountInfo fetches and returns info about the user's account.
func (c *Client) GetAccountInfo(ctx context.Context) (AccountInfo, error) {
	c.logger.Debug("Getting account info...", zapDebridService)

	resBytes, err := c.get(ctx, c.opts.BaseURL+"/account/info", nil)
	if err != nil {
		return AccountInfo{}, fmt.Errorf("couldn't get account info: %w", err)
	}
	if gjson.GetBytes(resBytes, "status").String() != "success" {
		message := gjson.GetBytes(resBytes, "message").String()
		return AccountInfo{}, fmt.Errorf("got error response from Premiumize: %v", message)
	}
	accInfo := AccountInfo{}
	if err = json.Unmarshal(resBytes, &accInfo); err != nil {
		return AccountInfo{}, fmt.Errorf("couldn't unmarshal account info: %w", err)
	}

	c.logger.Debug("Got account info", zap.String("accountInfo", fmt.Sprintf("%+v", accInfo)), zapDebridService)
	return accInfo, nil
}

// CheckCache checks if files are already in Premiumize's cache.
// An item can be any link that Premiumize supports: Containers, direct links, magnet URLs, torrent info hashes.
// The returned map contains only entries for cached files and uses the item as key.
func (c *Client) CheckCache(ctx context.Context, items ...string) (map[string]CachedFile, error) {
	c.logger.Debug("Checking cache...", zapDebridService)

	data := url.Values{"items[]": items}
	resBytes, err := c.get(ctx, c.opts.BaseURL+"/cache/check", data)
	if err != nil {
		return nil, fmt.Errorf("couldn't check cache: %w", err)
	}
	if gjson.GetBytes(resBytes, "status").String() != "success" {
		message := gjson.GetBytes(resBytes, "message").String()
		return nil, fmt.Errorf("got error response from Premiumize: %v", message)
	}
	cachedFiles := make(map[string]CachedFile, len(items))
	for i, item := range items {
		if gjson.GetBytes(resBytes, "response."+strconv.Itoa(i)).Bool() {
			cachedFiles[item] = CachedFile{
				Transcoded: gjson.GetBytes(resBytes, "transcoded."+strconv.Itoa(i)).Bool(),
				Filename:   gjson.GetBytes(resBytes, "filename."+strconv.Itoa(i)).String(),
				Filesize:   gjson.GetBytes(resBytes, "filesize."+strconv.Itoa(i)).String(),
			}
		}
	}

	c.logger.Debug("Checked cache", zap.String("cachedFiles", fmt.Sprintf("%+v", cachedFiles)), zapDebridService)
	return cachedFiles, nil
}
