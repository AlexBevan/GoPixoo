package pixoo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client communicates with a Pixoo64 device over HTTP.
type Client struct {
	IP         string
	HTTPClient *http.Client
}

// NewClient creates a Client for the given device IP.
func NewClient(ip string) *Client {
	return &Client{
		IP: ip,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) url() string {
	return fmt.Sprintf("http://%s:80/post", c.IP)
}

// Post sends a JSON command to the device and returns the parsed response.
func (c *Client) Post(payload map[string]interface{}) (map[string]interface{}, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	resp, err := c.HTTPClient.Post(c.url(), "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("post to device: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return result, nil
}
