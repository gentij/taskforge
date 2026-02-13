package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type Envelope struct {
	Ok         bool            `json:"ok"`
	StatusCode int             `json:"statusCode"`
	Path       string          `json:"path"`
	Timestamp  string          `json:"timestamp"`
	Data       json.RawMessage `json:"data"`
	Error      *APIError       `json:"error,omitempty"`
}

func NewClient(baseURL string, token string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) GetJSON(path string, out any) error {
	return c.doJSON(http.MethodGet, path, nil, out)
}

func (c *Client) PostJSON(path string, body any, out any) error {
	return c.doJSON(http.MethodPost, path, body, out)
}

func (c *Client) doJSON(method string, path string, body any, out any) error {
	fullURL, err := c.buildURL(path)
	if err != nil {
		return err
	}

	var bodyReader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(payload)
	}

	req, err := http.NewRequest(method, fullURL, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if strings.TrimSpace(c.Token) != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var env Envelope
	if err := json.Unmarshal(data, &env); err != nil {
		return err
	}

	if !env.Ok {
		if env.Error != nil {
			return fmt.Errorf("%s: %s", env.Error.Code, env.Error.Message)
		}
		return errors.New("request failed")
	}

	if out != nil && len(env.Data) > 0 {
		if err := json.Unmarshal(env.Data, out); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) buildURL(path string) (string, error) {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path, nil
	}

	base, err := url.Parse(c.BaseURL + "/")
	if err != nil {
		return "", err
	}

	ref, err := url.Parse(strings.TrimLeft(path, "/"))
	if err != nil {
		return "", err
	}

	return base.ResolveReference(ref).String(), nil
}
