package pixoo

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// ParseHexColor parses a hex color string like "#FF0000" or "FF0000" into RGB components.
func ParseHexColor(hex string) (uint8, uint8, uint8, error) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return 0, 0, 0, fmt.Errorf("invalid hex color %q: must be 6 hex digits", hex)
	}
	var r, g, b uint8
	_, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid hex color %q: %w", hex, err)
	}
	return r, g, b, nil
}

// BuildSolidFrame returns base64-encoded RGB pixel data for a solid color frame.
func BuildSolidFrame(size int, r, g, b uint8) string {
	data := make([]byte, size*size*3)
	for i := 0; i < size*size; i++ {
		data[i*3] = r
		data[i*3+1] = g
		data[i*3+2] = b
	}
	return base64.StdEncoding.EncodeToString(data)
}

// BuildPixelFrame returns base64-encoded RGB pixel data with a single pixel set (rest black).
func BuildPixelFrame(size int, x, y int, r, g, b uint8) string {
	data := make([]byte, size*size*3)
	offset := (y*size + x) * 3
	data[offset] = r
	data[offset+1] = g
	data[offset+2] = b
	return base64.StdEncoding.EncodeToString(data)
}
