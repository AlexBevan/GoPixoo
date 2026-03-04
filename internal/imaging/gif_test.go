package imaging

// Tests for GIF frame extraction via ExtractGIFFrames.

import (
	"image"
	"image/color"
	"image/color/palette"
	"image/gif"
	"image/png"
	"os"
	"testing"
)

// makeTestGIF creates a temporary GIF file with the given number of frames.
// Each frame is width x height with a distinct solid color.
// Returns the path to the temp file. Caller must remove it.
func makeTestGIF(t *testing.T, width, height, frames int, delayMs int) string {
	t.Helper()

	g := &gif.GIF{}
	for i := 0; i < frames; i++ {
		img := image.NewPaletted(image.Rect(0, 0, width, height), palette.Plan9)
		// Fill each frame with a different color
		c := palette.Plan9[i%len(palette.Plan9)]
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				img.Set(x, y, c)
			}
		}
		g.Image = append(g.Image, img)
		g.Delay = append(g.Delay, delayMs/10) // GIF delays are in 10ms units
	}

	f, err := os.CreateTemp("", "testgif-*.gif")
	if err != nil {
		t.Fatal(err)
	}
	if err := gif.EncodeAll(f, g); err != nil {
		f.Close()
		os.Remove(f.Name())
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

// makeTestPNG creates a temporary PNG file.
func makeTestPNG(t *testing.T, width, height int) string {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 100, G: 150, B: 200, A: 255})
		}
	}

	f, err := os.CreateTemp("", "testpng-*.png")
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		os.Remove(f.Name())
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestExtractGIFFrames_MultiFrameGIF(t *testing.T) {
	path := makeTestGIF(t, 4, 4, 3, 100)
	defer os.Remove(path)

	frames, delays, err := ExtractGIFFrames(path, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(frames) != 3 {
		t.Errorf("expected 3 frames, got %d", len(frames))
	}
	if len(delays) != 3 {
		t.Errorf("expected 3 delays, got %d", len(delays))
	}

	// Each frame should be 4x4
	for i, frame := range frames {
		b := frame.Bounds()
		if b.Dx() != 4 || b.Dy() != 4 {
			t.Errorf("frame %d: expected 4x4, got %dx%d", i, b.Dx(), b.Dy())
		}
	}
}

func TestExtractGIFFrames_DelaysCorrect(t *testing.T) {
	path := makeTestGIF(t, 4, 4, 2, 200)
	defer os.Remove(path)

	_, delays, err := ExtractGIFFrames(path, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i, d := range delays {
		// Delays should be in milliseconds; we set 200ms
		if d != 200 {
			t.Errorf("frame %d delay: expected 200ms, got %dms", i, d)
		}
	}
}

func TestExtractGIFFrames_SingleFrameGIF(t *testing.T) {
	path := makeTestGIF(t, 4, 4, 1, 0)
	defer os.Remove(path)

	frames, delays, err := ExtractGIFFrames(path, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(frames) != 1 {
		t.Errorf("expected 1 frame, got %d", len(frames))
	}
	if len(delays) != 1 {
		t.Errorf("expected 1 delay, got %d", len(delays))
	}
}

func TestExtractGIFFrames_PNGReturnsSingleFrame(t *testing.T) {
	path := makeTestPNG(t, 8, 8)
	defer os.Remove(path)

	frames, delays, err := ExtractGIFFrames(path, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(frames) != 1 {
		t.Errorf("expected 1 frame for PNG, got %d", len(frames))
	}
	if len(delays) != 1 {
		t.Errorf("expected 1 delay entry for PNG, got %d", len(delays))
	}
}

func TestExtractGIFFrames_MissingFile(t *testing.T) {
	_, _, err := ExtractGIFFrames("/nonexistent/path/to/animation.gif", 100)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestExtractGIFFrames_FrameColorsDistinct(t *testing.T) {
	path := makeTestGIF(t, 4, 4, 2, 100)
	defer os.Remove(path)

	frames, _, err := ExtractGIFFrames(path, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(frames) < 2 {
		t.Skip("need at least 2 frames")
	}

	// Frames should have different content (different fill colors)
	r0, g0, b0, _ := frames[0].At(0, 0).RGBA()
	r1, g1, b1, _ := frames[1].At(0, 0).RGBA()
	if r0 == r1 && g0 == g1 && b0 == b1 {
		t.Error("frames should have distinct colors but pixel (0,0) is identical")
	}
}
