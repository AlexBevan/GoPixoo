package pixoo

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestPost_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("invalid JSON from client: %v", err)
		}
		if payload["Command"] != "Draw/SendHttpGif" {
			t.Errorf("unexpected command: %v", payload["Command"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error_code": 0,
		})
	}))
	defer server.Close()

	// Extract host (without scheme) from test server URL
	addr := strings.TrimPrefix(server.URL, "http://")
	host := strings.Split(addr, ":")[0]

	c := &Client{
		IP:         host,
		HTTPClient: server.Client(),
	}
	// Override URL to use test server port
	origURL := c.url
	_ = origURL
	c.HTTPClient.Transport = &http.Transport{}
	// We need to redirect to the test server, so replace the client's HTTP client
	c.HTTPClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			req.URL.Scheme = "http"
			req.URL.Host = strings.TrimPrefix(server.URL, "http://")
			return http.DefaultTransport.RoundTrip(req)
		}),
	}

	payload := map[string]interface{}{
		"Command": CmdSendHTTPGif,
		"PicData": "AAAA",
	}
	result, err := c.Post(payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["error_code"] != float64(0) {
		t.Errorf("expected error_code 0, got %v", result["error_code"])
	}
}

// roundTripFunc allows using a function as an http.RoundTripper.
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestPost_ConnectionError(t *testing.T) {
	c := NewClient("192.0.2.1") // RFC 5737 TEST-NET, guaranteed unreachable
	c.HTTPClient.Timeout = 500 * time.Millisecond

	_, err := c.Post(map[string]interface{}{"Command": "test"})
	if err == nil {
		t.Fatal("expected error for unreachable host, got nil")
	}
	if !strings.Contains(err.Error(), "post to device") {
		t.Errorf("expected wrapped error, got: %v", err)
	}
}

func TestPost_MalformedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("this is not json"))
	}))
	defer server.Close()

	c := &Client{
		IP: "127.0.0.1",
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				req.URL.Scheme = "http"
				req.URL.Host = strings.TrimPrefix(server.URL, "http://")
				return http.DefaultTransport.RoundTrip(req)
			}),
		},
	}

	_, err := c.Post(map[string]interface{}{"Command": "test"})
	if err == nil {
		t.Fatal("expected error for malformed response, got nil")
	}
	if !strings.Contains(err.Error(), "unmarshal response") {
		t.Errorf("expected unmarshal error, got: %v", err)
	}
}

func TestPost_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		json.NewEncoder(w).Encode(map[string]interface{}{"error_code": 0})
	}))
	defer server.Close()

	c := &Client{
		IP: "127.0.0.1",
		HTTPClient: &http.Client{
			Timeout: 200 * time.Millisecond,
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				req.URL.Scheme = "http"
				req.URL.Host = strings.TrimPrefix(server.URL, "http://")
				return http.DefaultTransport.RoundTrip(req)
			}),
		},
	}

	_, err := c.Post(map[string]interface{}{"Command": "test"})
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestPost_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// empty body
	}))
	defer server.Close()

	c := &Client{
		IP: "127.0.0.1",
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				req.URL.Scheme = "http"
				req.URL.Host = strings.TrimPrefix(server.URL, "http://")
				return http.DefaultTransport.RoundTrip(req)
			}),
		},
	}

	_, err := c.Post(map[string]interface{}{"Command": "test"})
	if err == nil {
		t.Fatal("expected error for empty response body, got nil")
	}
}

func TestNewClient_Defaults(t *testing.T) {
	c := NewClient("10.0.0.1")
	if c.IP != "10.0.0.1" {
		t.Errorf("expected IP 10.0.0.1, got %s", c.IP)
	}
	if c.HTTPClient == nil {
		t.Fatal("expected non-nil HTTPClient")
	}
	if c.HTTPClient.Timeout != 10*time.Second {
		t.Errorf("expected 10s timeout, got %v", c.HTTPClient.Timeout)
	}
}

func TestClient_URL(t *testing.T) {
	c := NewClient("192.168.1.100")
	expected := "http://192.168.1.100:80/post"
	if got := c.url(); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
