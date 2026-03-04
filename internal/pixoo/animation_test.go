package pixoo

import (
	"testing"
)

func TestSendGIF_Structure(t *testing.T) {
	payload := SendGIF(5, 2, 42, 100, 64, "base64data==")

	tests := []struct {
		key      string
		expected interface{}
	}{
		{"Command", CmdSendHTTPGif},
		{"PicNum", 5},
		{"PicWidth", 64},
		{"PicOffset", 2},
		{"PicID", 42},
		{"PicSpeed", 100},
		{"PicData", "base64data=="},
	}

	for _, tc := range tests {
		val, ok := payload[tc.key]
		if !ok {
			t.Errorf("missing key %q in SendGIF payload", tc.key)
			continue
		}
		if val != tc.expected {
			t.Errorf("key %q: expected %v (%T), got %v (%T)", tc.key, tc.expected, tc.expected, val, val)
		}
	}

	if len(payload) != 7 {
		t.Errorf("expected 7 fields in SendGIF payload, got %d", len(payload))
	}
}

func TestSendGIF_CommandName(t *testing.T) {
	payload := SendGIF(1, 0, 1, 50, 64, "")
	if payload["Command"] != "Draw/SendHttpGif" {
		t.Errorf("expected Draw/SendHttpGif, got %v", payload["Command"])
	}
}

func TestSendGIF_ZeroValues(t *testing.T) {
	payload := SendGIF(0, 0, 0, 0, 0, "")
	if payload["PicNum"] != 0 {
		t.Errorf("expected PicNum 0, got %v", payload["PicNum"])
	}
	if payload["PicSpeed"] != 0 {
		t.Errorf("expected PicSpeed 0, got %v", payload["PicSpeed"])
	}
	if payload["PicData"] != "" {
		t.Errorf("expected empty PicData, got %v", payload["PicData"])
	}
}

func TestResetGIFID_Structure(t *testing.T) {
	payload := ResetGIFID()

	cmd, ok := payload["Command"]
	if !ok {
		t.Fatal("missing Command key in ResetGIFID payload")
	}
	if cmd != CmdResetHTTPGifID {
		t.Errorf("expected %q, got %v", CmdResetHTTPGifID, cmd)
	}
	if cmd != "Draw/ResetHttpGifId" {
		t.Errorf("expected literal Draw/ResetHttpGifId, got %v", cmd)
	}
	if len(payload) != 1 {
		t.Errorf("expected 1 field in ResetGIFID payload, got %d", len(payload))
	}
}

func TestGetGIFID_Structure(t *testing.T) {
	payload := GetGIFID()

	cmd, ok := payload["Command"]
	if !ok {
		t.Fatal("missing Command key in GetGIFID payload")
	}
	if cmd != CmdGetHTTPGifID {
		t.Errorf("expected %q, got %v", CmdGetHTTPGifID, cmd)
	}
	if cmd != "Draw/GetHttpGifId" {
		t.Errorf("expected literal Draw/GetHttpGifId, got %v", cmd)
	}
	if len(payload) != 1 {
		t.Errorf("expected 1 field in GetGIFID payload, got %d", len(payload))
	}
}

func TestCommandConstants(t *testing.T) {
	// Verify the command constants match the Pixoo64 API exactly
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"SendHTTPGif", CmdSendHTTPGif, "Draw/SendHttpGif"},
		{"ResetHTTPGifID", CmdResetHTTPGifID, "Draw/ResetHttpGifId"},
		{"GetHTTPGifID", CmdGetHTTPGifID, "Draw/GetHttpGifId"},
	}
	for _, tc := range tests {
		if tc.constant != tc.expected {
			t.Errorf("%s: expected %q, got %q", tc.name, tc.expected, tc.constant)
		}
	}
}

func TestSendGIF_PicWidthAlways64(t *testing.T) {
	// PicWidth should always be 64 regardless of input params
	for _, picNum := range []int{1, 5, 10, 100} {
		payload := SendGIF(picNum, 0, 1, 50, 64, "data")
		if payload["PicWidth"] != 64 {
			t.Errorf("PicWidth should be 64 when passed 64, got %v for picNum=%d", payload["PicWidth"], picNum)
		}
	}
}
