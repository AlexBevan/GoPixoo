package imaging

import (
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"os"
	"strings"
)

// ExtractGIFFrames loads a GIF and returns composited frames with per-frame delays in ms.
// If the file is not a GIF, it falls back to loading a single static image with the given default delay.
func ExtractGIFFrames(path string, defaultDelayMs int) ([]image.Image, []int, error) {
	if !isGIF(path) {
		img, err := LoadImage(path)
		if err != nil {
			return nil, nil, err
		}
		return []image.Image{img}, []int{defaultDelayMs}, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("open gif: %w", err)
	}
	defer f.Close()

	g, err := gif.DecodeAll(f)
	if err != nil {
		return nil, nil, fmt.Errorf("decode gif: %w", err)
	}

	if len(g.Image) == 0 {
		return nil, nil, fmt.Errorf("gif contains no frames")
	}

	// Build a canvas matching the GIF's logical screen size.
	canvasWidth := g.Config.Width
	canvasHeight := g.Config.Height
	if canvasWidth == 0 || canvasHeight == 0 {
		canvasWidth = g.Image[0].Bounds().Dx()
		canvasHeight = g.Image[0].Bounds().Dy()
	}

	canvas := image.NewRGBA(image.Rect(0, 0, canvasWidth, canvasHeight))
	frames := make([]image.Image, 0, len(g.Image))
	delays := make([]int, 0, len(g.Image))

	for i, frame := range g.Image {
		// Draw frame onto canvas respecting disposal method.
		draw.Draw(canvas, frame.Bounds(), frame, frame.Bounds().Min, draw.Over)

		// Snapshot the composited canvas.
		snapshot := image.NewRGBA(canvas.Bounds())
		draw.Draw(snapshot, snapshot.Bounds(), canvas, image.Point{}, draw.Src)
		frames = append(frames, snapshot)

		// GIF delays are in 100ths of a second; convert to ms.
		delayMs := defaultDelayMs
		if i < len(g.Delay) && g.Delay[i] > 0 {
			delayMs = g.Delay[i] * 10
		}
		delays = append(delays, delayMs)

		// Handle disposal method for next frame.
		disposal := byte(gif.DisposalNone)
		if i < len(g.Disposal) {
			disposal = g.Disposal[i]
		}
		switch disposal {
		case gif.DisposalBackground:
			draw.Draw(canvas, frame.Bounds(), image.Transparent, image.Point{}, draw.Src)
		case gif.DisposalPrevious:
			// Restore canvas to state before this frame was drawn.
			// For simplicity, clear the frame region (close approximation).
			draw.Draw(canvas, frame.Bounds(), image.Transparent, image.Point{}, draw.Src)
		}
	}

	return frames, delays, nil
}

func isGIF(path string) bool {
	return strings.HasSuffix(strings.ToLower(path), ".gif")
}
