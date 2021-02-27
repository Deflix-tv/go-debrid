package premiumize

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"go.uber.org/zap"
)

func (c *Client) get(ctx context.Context, urlString string, data url.Values) ([]byte, error) {
	if c.auth.OAuth2 {
		urlString += "?access_token=" + c.auth.KeyOrToken
	} else {
		urlString += "?apikey=" + c.auth.KeyOrToken
	}

	// map[string][]string
	for k, vals := range data {
		for _, val := range vals {
			urlString += "&" + url.QueryEscape(k) + "=" + url.QueryEscape(val)
		}
	}

	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't create GET request: %w", err)
	}
	for headerKey, headerVal := range c.opts.ExtraHeaders {
		req.Header.Add(headerKey, headerVal)
	}

	c.logger.Debug("Sending request", zap.String("request", fmt.Sprintf("%+v", req)), zapDebridService)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("couldn't send GET request: %w", err)
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	c.logger.Debug("Got response", zap.Int("status", res.StatusCode), zap.NamedError("bodyReadError", err), zap.ByteString("response", resBody), zapDebridService)

	// Check server response status
	if res.StatusCode != http.StatusOK {
		// resBody can be nil if above ioutil.ReadAll failed, but in that case we don't care about the related error.
		return resBody, fmt.Errorf("bad HTTP response status: %v", res.Status)
	}

	if err != nil {
		return nil, fmt.Errorf("couldn't read response body: %w", err)
	}
	return resBody, nil
}

func (c *Client) post(ctx context.Context, urlString string, data url.Values, form bool) ([]byte, error) {
	if c.auth.OAuth2 {
		urlString += "?access_token=" + c.auth.KeyOrToken
	} else {
		urlString += "?apikey=" + c.auth.KeyOrToken
	}

	var req *http.Request
	var err error
	if form {
		req, err = http.NewRequest("POST", urlString, strings.NewReader(data.Encode()))
	} else {
		// map[string][]string
		for k, vals := range data {
			for _, val := range vals {
				urlString += "&" + url.QueryEscape(k) + "=" + url.QueryEscape(val)
			}
		}
		req, err = http.NewRequest("POST", urlString, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("couldn't create POST request: %w", err)
	}
	if form {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	for headerKey, headerVal := range c.opts.ExtraHeaders {
		req.Header.Add(headerKey, headerVal)
	}

	c.logger.Debug("Sending request", zap.String("request", fmt.Sprintf("%+v", req)), zapDebridService)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("couldn't send POST request: %w", err)
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	c.logger.Debug("Got response", zap.Int("status", res.StatusCode), zap.NamedError("bodyReadError", err), zap.ByteString("response", resBody), zapDebridService)

	// Check server response.
	if res.StatusCode != http.StatusOK {
		// resBody can be nil if above ioutil.ReadAll failed, but in that case we don't care about the related error.
		return resBody, fmt.Errorf("bad HTTP response status: %v", res.Status)
	}

	if err != nil {
		return nil, fmt.Errorf("couldn't read response body: %w", err)
	}
	return resBody, nil
}
