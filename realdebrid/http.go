package realdebrid

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"go.uber.org/zap"
)

// data can be nil
func (c *Client) get(ctx context.Context, url string, data url.Values) ([]byte, error) {
	var body io.Reader
	if len(data) > 0 {
		body = strings.NewReader(data.Encode())
	}
	req, err := http.NewRequest("GET", url, body)
	if err != nil {
		return nil, fmt.Errorf("couldn't create GET request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.auth.KeyOrToken)
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
		if err, found := errMap[res.StatusCode]; found {
			return resBody, err
		}
		// resBody can be nil if above ioutil.ReadAll failed, but in that case we don't care about the related error.
		return resBody, fmt.Errorf("bad HTTP response status: %v", res.Status)
	}

	if err != nil {
		return nil, fmt.Errorf("couldn't read response body: %w", err)
	}
	return resBody, nil
}

func (c *Client) post(ctx context.Context, url string, data url.Values) ([]byte, error) {
	// RealDebrid asks for the original IP for all POST requests.
	if c.opts.ForwardOriginIP {
		if c.auth.IP == "" {
			return nil, errors.New("auth.IP is empty but client is configured to forward the user's original IP")
		}
		data.Add("ip", c.auth.IP)
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("couldn't create POST request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.auth.KeyOrToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
	// Different RealDebrid API POST endpoints return different status codes.
	if res.StatusCode != http.StatusCreated &&
		res.StatusCode != http.StatusNoContent &&
		res.StatusCode != http.StatusOK {
		if err, found := errMap[res.StatusCode]; found {
			return resBody, err
		}
		// resBody can be nil if above ioutil.ReadAll failed, but in that case we don't care about the related error.
		return resBody, fmt.Errorf("bad HTTP response status: %v", res.Status)
	}

	if err != nil {
		return nil, fmt.Errorf("couldn't read response body: %w", err)
	}
	return resBody, nil
}
