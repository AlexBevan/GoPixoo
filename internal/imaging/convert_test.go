package imaging

import (
	"encoding/base64"
	"image"
	"image/color"
	"os"
	"testing"
)

func TestEncodePixels_2x2KnownImage(t *testing.T) {
	// Create a 2x2 image with known pixel colors
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})     // red
	img.Set(1, 0, color.RGBA{R: 0, G: 255, B: 0, A: 255})     // green
	img.Set(0, 1, color.RGBA{R: 0, G: 0, B: 255, A: 255})     // blue
	img.Set(1, 1, color.RGBA{R: 255, G: 255, B: 255, A: 255}) // white

	encoded := EncodePixels(img)

	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("invalid base64: %v", err)
	}

	// Output should be exactly 64*64*3 = 12288 bytes
	expectedLen := PixooSize * PixooSize * 3
	if len(data) != expectedLen {
		t.Fatalf("expected %d bytes, got %d", expectedLen, len(data))
	}

	// Check first row pixels: (0,0)=red, (1,0)=green, (2..63,0)=black
	assertPixel(t, data, 0, 0, 255, 0, 0, "top-left red")
	assertPixel(t, data, 1, 0, 0, 255, 0, "top-right green")
	assertPixel(t, data, 2, 0, 0, 0, 0, "padding pixel (2,0)")

	// Check second row: (0,1)=blue, (1,1)=white
	assertPixel(t, data, 0, 1, 0, 0, 255, "bottom-left blue")
	assertPixel(t, data, 1, 1, 255, 255, 255, "bottom-right white")

	// Check a padding pixel deep in the image
	assertPixel(t, data, 63, 63, 0, 0, 0, "padding pixel (63,63)")
}

func TestEncodePixels_SmallImagePadsWithBlack(t *testing.T) {
	// 1x1 red image — everything else should be black
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{R: 128, G: 64, B: 32, A: 255})

	encoded := EncodePixels(img)
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("invalid base64: %v", err)
	}

	// First pixel should be our color
	assertPixel(t, data, 0, 0, 128, 64, 32, "origin pixel")

	// Spot check several padding pixels
	for _, pt := range [][2]int{{1, 0}, {0, 1}, {32, 32}, {63, 63}} {
		assertPixel(t, data, pt[0], pt[1], 0, 0, 0, "padding")
	}
}

func TestEncodePixels_64x64Image(t *testing.T) {
	// Full 64x64 image filled with a single color
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			img.Set(x, y, color.RGBA{R: 10, G: 20, B: 30, A: 255})
		}
	}

	encoded := EncodePixels(img)
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("invalid base64: %v", err)
	}

	// Every pixel should be (10, 20, 30)
	for i := 0; i < len(data); i += 3 {
		if data[i] != 10 || data[i+1] != 20 || data[i+2] != 30 {
			px := i / 3
			t.Errorf("pixel %d: expected (10,20,30), got (%d,%d,%d)", px, data[i], data[i+1], data[i+2])
			break
		}
	}
}

func TestEncodePixels_OversizedImageClamps(t *testing.T) {
	// Image larger than 64x64 — should only read the top-left 64x64 region
	img := image.NewRGBA(image.Rect(0, 0, 128, 128))
	for y := 0; y < 128; y++ {
		for x := 0; x < 128; x++ {
			if x < 64 && y < 64 {
				img.Set(x, y, color.RGBA{R: 100, G: 100, B: 100, A: 255})
			} else {
				img.Set(x, y, color.RGBA{R: 200, G: 200, B: 200, A: 255})
			}
		}
	}

	encoded := EncodePixels(img)
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("invalid base64: %v", err)
	}

	if len(data) != PixooSize*PixooSize*3 {
		t.Fatalf("expected %d bytes, got %d", PixooSize*PixooSize*3, len(data))
	}

	// All pixels should be (100,100,100) since we only read the 64x64 region
	for i := 0; i < len(data); i += 3 {
		if data[i] != 100 || data[i+1] != 100 || data[i+2] != 100 {
			px := i / 3
			t.Errorf("pixel %d: expected (100,100,100), got (%d,%d,%d)", px, data[i], data[i+1], data[i+2])
			break
		}
	}
}

func TestEncodePixels_OutputLength(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	encoded := EncodePixels(img)
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("invalid base64: %v", err)
	}
	expected := PixooSize * PixooSize * 3
	if len(data) != expected {
		t.Errorf("expected %d raw bytes, got %d", expected, len(data))
	}
}

func TestLoadImage_NonexistentFile(t *testing.T) {
	_, err := LoadImage("/nonexistent/path/to/image.png")
	if err == nil {
		t.Fatal("expected error for nonexistent file, got nil")
	}
}

func TestLoadImage_InvalidImageData(t *testing.T) {
	// Create a temp file with non-image content
	f, err := os.CreateTemp("", "notanimage-*.png")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	f.WriteString("this is not an image")
	f.Close()

	_, err = LoadImage(f.Name())
	if err == nil {
		t.Fatal("expected error for invalid image data, got nil")
	}
}

func TestPixooSizeConstant(t *testing.T) {
	if PixooSize != 64 {
		t.Errorf("PixooSize should be 64, got %d", PixooSize)
	}
}

func TestEncodePixels_NonZeroOriginBounds(t *testing.T) {
	// Image with non-zero origin (e.g. sub-image)
	parent := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			parent.Set(x, y, color.RGBA{R: uint8(x * 25), G: uint8(y * 25), B: 0, A: 255})
		}
	}
	sub := parent.SubImage(image.Rect(5, 5, 8, 8))

	encoded := EncodePixels(sub)
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("invalid base64: %v", err)
	}

	// Pixel (0,0) of the sub-image corresponds to parent pixel (5,5)
	assertPixel(t, data, 0, 0, 125, 125, 0, "sub-image origin")
}

// assertPixel checks a pixel at (x,y) in the raw RGB data buffer.
func assertPixel(t *testing.T, data []byte, x, y int, r, g, b byte, label string) {
	t.Helper()
	idx := (y*PixooSize + x) * 3
	if idx+2 >= len(data) {
		t.Errorf("%s: pixel (%d,%d) index out of range", label, x, y)
		return
	}
	if data[idx] != r || data[idx+1] != g || data[idx+2] != b {
		t.Errorf("%s: pixel (%d,%d) expected (%d,%d,%d), got (%d,%d,%d)",
			label, x, y, r, g, b, data[idx], data[idx+1], data[idx+2])
	}
}
