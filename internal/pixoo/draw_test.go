package pixoo

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestParseHexColor_ValidColors(t *testing.T) {
	tests := []struct {
		name  string
		hex   string
		wantR uint8
		wantG uint8
		wantB uint8
	}{
		{"red with hash", "#FF0000", 255, 0, 0},
		{"red without hash", "FF0000", 255, 0, 0},
		{"green", "#00FF00", 0, 255, 0},
		{"blue", "#0000FF", 0, 0, 255},
		{"white", "#FFFFFF", 255, 255, 255},
		{"black", "#000000", 0, 0, 0},
		{"cyan", "#00FFFF", 0, 255, 255},
		{"magenta", "#FF00FF", 255, 0, 255},
		{"yellow", "#FFFF00", 255, 255, 0},
		{"lowercase", "#aabbcc", 170, 187, 204},
		{"mixed case", "#AaBbCc", 170, 187, 204},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, g, b, err := ParseHexColor(tt.hex)
			if err != nil {
				t.Fatalf("ParseHexColor(%q) unexpected error: %v", tt.hex, err)
			}
			if r != tt.wantR || g != tt.wantG || b != tt.wantB {
				t.Errorf("ParseHexColor(%q) = (%d, %d, %d), want (%d, %d, %d)",
					tt.hex, r, g, b, tt.wantR, tt.wantG, tt.wantB)
			}
		})
	}
}

func TestParseHexColor_InvalidColors(t *testing.T) {
	tests := []struct {
		name string
		hex  string
	}{
		{"too short", "#FFF"},
		{"too long", "#FFFFFFF"},
		{"invalid chars", "#GGGGGG"},
		{"invalid chars mixed", "#FF00GG"},
		{"not hex at all", "not-a-color"},
		{"empty string", ""},
		{"just hash", "#"},
		{"five chars", "#FFFFF"},
		{"seven chars", "#FFFFFFF"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, err := ParseHexColor(tt.hex)
			if err == nil {
				t.Errorf("ParseHexColor(%q) expected error, got nil", tt.hex)
			}
		})
	}
}

func TestBuildSolidFrame_Structure(t *testing.T) {
	tests := []struct {
		name string
		size int
		r    uint8
		g    uint8
		b    uint8
	}{
		{"64x64 red", 64, 255, 0, 0},
		{"64x64 green", 64, 0, 255, 0},
		{"64x64 blue", 64, 0, 0, 255},
		{"64x64 white", 64, 255, 255, 255},
		{"64x64 black", 64, 0, 0, 0},
		{"16x16 cyan", 16, 0, 255, 255},
		{"32x32 magenta", 32, 255, 0, 255},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := BuildSolidFrame(tt.size, tt.r, tt.g, tt.b)

			// Verify it's valid base64
			decoded, err := base64.StdEncoding.DecodeString(encoded)
			if err != nil {
				t.Fatalf("BuildSolidFrame returned invalid base64: %v", err)
			}

			// Verify length
			expectedLen := tt.size * tt.size * 3
			if len(decoded) != expectedLen {
				t.Errorf("decoded length = %d, want %d", len(decoded), expectedLen)
			}

			// Spot-check first pixel
			if decoded[0] != tt.r || decoded[1] != tt.g || decoded[2] != tt.b {
				t.Errorf("first pixel = (%d, %d, %d), want (%d, %d, %d)",
					decoded[0], decoded[1], decoded[2], tt.r, tt.g, tt.b)
			}

			// Spot-check middle pixel
			mid := (tt.size*tt.size/2) * 3
			if decoded[mid] != tt.r || decoded[mid+1] != tt.g || decoded[mid+2] != tt.b {
				t.Errorf("middle pixel = (%d, %d, %d), want (%d, %d, %d)",
					decoded[mid], decoded[mid+1], decoded[mid+2], tt.r, tt.g, tt.b)
			}

			// Spot-check last pixel
			last := len(decoded) - 3
			if decoded[last] != tt.r || decoded[last+1] != tt.g || decoded[last+2] != tt.b {
				t.Errorf("last pixel = (%d, %d, %d), want (%d, %d, %d)",
					decoded[last], decoded[last+1], decoded[last+2], tt.r, tt.g, tt.b)
			}
		})
	}
}

func TestBuildPixelFrame_Structure(t *testing.T) {
	tests := []struct {
		name string
		size int
		x    int
		y    int
		r    uint8
		g    uint8
		b    uint8
	}{
		{"64x64 top-left red", 64, 0, 0, 255, 0, 0},
		{"64x64 center green", 64, 32, 32, 0, 255, 0},
		{"64x64 bottom-right blue", 64, 63, 63, 0, 0, 255},
		{"16x16 center white", 16, 8, 8, 255, 255, 255},
		{"32x32 corner cyan", 32, 0, 0, 0, 255, 255},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := BuildPixelFrame(tt.size, tt.x, tt.y, tt.r, tt.g, tt.b)

			// Verify it's valid base64
			decoded, err := base64.StdEncoding.DecodeString(encoded)
			if err != nil {
				t.Fatalf("BuildPixelFrame returned invalid base64: %v", err)
			}

			// Verify length
			expectedLen := tt.size * tt.size * 3
			if len(decoded) != expectedLen {
				t.Errorf("decoded length = %d, want %d", len(decoded), expectedLen)
			}

			// Calculate expected offset for the pixel
			offset := (tt.y*tt.size + tt.x) * 3

			// Check that the specified pixel has the right color
			if decoded[offset] != tt.r || decoded[offset+1] != tt.g || decoded[offset+2] != tt.b {
				t.Errorf("pixel at (%d, %d) = (%d, %d, %d), want (%d, %d, %d)",
					tt.x, tt.y, decoded[offset], decoded[offset+1], decoded[offset+2], tt.r, tt.g, tt.b)
			}

			// Check that another pixel is black (rest should be black)
			checkOtherX := (tt.x + 1) % tt.size
			checkOtherY := (tt.y + 1) % tt.size
			otherOffset := (checkOtherY*tt.size + checkOtherX) * 3

			// Only check if it's a different pixel
			if otherOffset != offset {
				if decoded[otherOffset] != 0 || decoded[otherOffset+1] != 0 || decoded[otherOffset+2] != 0 {
					t.Errorf("other pixel at (%d, %d) = (%d, %d, %d), want (0, 0, 0)",
						checkOtherX, checkOtherY, decoded[otherOffset], decoded[otherOffset+1], decoded[otherOffset+2])
				}
			}
		})
	}
}

func TestBuildSolidFrame_AllPixelsCorrect(t *testing.T) {
	size := 8
	r, g, b := uint8(100), uint8(150), uint8(200)
	encoded := BuildSolidFrame(size, r, g, b)
	decoded, _ := base64.StdEncoding.DecodeString(encoded)

	// Verify every pixel is correct
	for i := 0; i < size*size; i++ {
		offset := i * 3
		if decoded[offset] != r || decoded[offset+1] != g || decoded[offset+2] != b {
			t.Errorf("pixel %d = (%d, %d, %d), want (%d, %d, %d)",
				i, decoded[offset], decoded[offset+1], decoded[offset+2], r, g, b)
		}
	}
}

func TestBuildPixelFrame_AllOtherPixelsBlack(t *testing.T) {
	size := 8
	targetX, targetY := 3, 4
	r, g, b := uint8(255), uint8(128), uint8(64)
	encoded := BuildPixelFrame(size, targetX, targetY, r, g, b)
	decoded, _ := base64.StdEncoding.DecodeString(encoded)

	targetOffset := (targetY*size + targetX) * 3

	// Verify every pixel
	for i := 0; i < size*size; i++ {
		offset := i * 3
		if offset == targetOffset {
			// This is our target pixel - should be colored
			if decoded[offset] != r || decoded[offset+1] != g || decoded[offset+2] != b {
				t.Errorf("target pixel %d = (%d, %d, %d), want (%d, %d, %d)",
					i, decoded[offset], decoded[offset+1], decoded[offset+2], r, g, b)
			}
		} else {
			// All other pixels should be black
			if decoded[offset] != 0 || decoded[offset+1] != 0 || decoded[offset+2] != 0 {
				t.Errorf("pixel %d = (%d, %d, %d), want (0, 0, 0)",
					i, decoded[offset], decoded[offset+1], decoded[offset+2])
			}
		}
	}
}

func TestParseHexColor_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		hex     string
		wantErr bool
	}{
		{"leading spaces", "  FF0000", true},
		{"trailing spaces", "FF0000  ", true},
		{"spaces in middle", "FF 00 00", true},
		{"multiple hashes", "##FF0000", true},
		{"hash in middle", "FF#0000", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, err := ParseHexColor(tt.hex)
			gotErr := err != nil
			if gotErr != tt.wantErr {
				t.Errorf("ParseHexColor(%q) error = %v, wantErr %v", tt.hex, err, tt.wantErr)
			}
		})
	}
}

func TestParseHexColor_ErrorMessages(t *testing.T) {
	tests := []struct {
		hex         string
		wantContain string
	}{
		{"#FFF", "must be 6 hex digits"},
		{"#FFFFFFF", "must be 6 hex digits"},
		{"#GGGGGG", "invalid hex color"},
	}

	for _, tt := range tests {
		t.Run(tt.hex, func(t *testing.T) {
			_, _, _, err := ParseHexColor(tt.hex)
			if err == nil {
				t.Fatalf("ParseHexColor(%q) expected error, got nil", tt.hex)
			}
			if !strings.Contains(err.Error(), tt.wantContain) {
				t.Errorf("error message %q does not contain %q", err.Error(), tt.wantContain)
			}
		})
	}
}
