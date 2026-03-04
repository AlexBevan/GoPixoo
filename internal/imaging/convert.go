package imaging

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	xdraw "golang.org/x/image/draw"
)

const PixooSize = 64

// LoadImage reads an image file from disk.
func LoadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open image: %w", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}
	return img, nil
}

// EncodePixels converts a 64x64 image to Pixoo64 base64 RGB format.
// Returns base64-encoded string of raw RGB pixel data (64*64*3 bytes).
func EncodePixels(img image.Image) string {
	return EncodePixelsSized(img, PixooSize)
}

// EncodePixelsSized converts an image to base64 RGB format for the given pixel size.
// Returns base64-encoded string of raw RGB pixel data (size*size*3 bytes).
func EncodePixelsSized(img image.Image, size int) string {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	if w > size {
		w = size
	}
	if h > size {
		h = size
	}

	data := make([]byte, 0, size*size*3)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if x < w && y < h {
				r, g, b, _ := img.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
				data = append(data, byte(r>>8), byte(g>>8), byte(b>>8))
			} else {
				data = append(data, 0, 0, 0)
			}
		}
	}
	return base64.StdEncoding.EncodeToString(data)
}

// ResizeMode controls how images are resized to the target dimensions.
type ResizeMode string

const (
	ResizeFit     ResizeMode = "fit"     // maintain aspect ratio, pad with black
	ResizeFill    ResizeMode = "fill"    // maintain aspect ratio, crop to anchor
	ResizeStretch ResizeMode = "stretch" // distort to fill target size
	ResizeNone    ResizeMode = "none"    // no resize
)

// CropAnchor controls where the crop is applied in fill mode.
type CropAnchor string

const (
	AnchorCenter CropAnchor = "center"
	AnchorTop    CropAnchor = "top"
	AnchorBottom CropAnchor = "bottom"
	AnchorLeft   CropAnchor = "left"
	AnchorRight  CropAnchor = "right"
)

// Resize scales an image to size×size pixels using the specified mode.
// Fill mode crops from center; use ResizeWithAnchor for other anchor positions.
func Resize(img image.Image, size int, mode ResizeMode) image.Image {
	return ResizeWithAnchor(img, size, mode, AnchorCenter)
}

// ResizeWithAnchor scales an image to size×size pixels using the specified mode
// and anchor. The anchor controls where the crop occurs in fill mode.
func ResizeWithAnchor(img image.Image, size int, mode ResizeMode, anchor CropAnchor) image.Image {
	if mode == ResizeNone {
		return img
	}

	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(dst, dst.Bounds(), image.NewUniform(color.Black), image.Point{}, draw.Src)

	srcBounds := img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	switch mode {
	case ResizeStretch:
		xdraw.CatmullRom.Scale(dst, dst.Bounds(), img, srcBounds, xdraw.Over, nil)

	case ResizeFit:
		scale := float64(size) / float64(srcW)
		if s := float64(size) / float64(srcH); s < scale {
			scale = s
		}
		newW := int(float64(srcW) * scale)
		newH := int(float64(srcH) * scale)
		offsetX := (size - newW) / 2
		offsetY := (size - newH) / 2
		dstRect := image.Rect(offsetX, offsetY, offsetX+newW, offsetY+newH)
		xdraw.CatmullRom.Scale(dst, dstRect, img, srcBounds, xdraw.Over, nil)

	case ResizeFill:
		scale := float64(size) / float64(srcW)
		if s := float64(size) / float64(srcH); s > scale {
			scale = s
		}
		newW := int(float64(srcW) * scale)
		newH := int(float64(srcH) * scale)

		offsetX, offsetY := fillOffset(newW, newH, size, anchor)

		tmp := image.NewRGBA(image.Rect(0, 0, newW, newH))
		xdraw.CatmullRom.Scale(tmp, tmp.Bounds(), img, srcBounds, xdraw.Over, nil)
		draw.Draw(dst, dst.Bounds(), tmp, image.Pt(offsetX, offsetY), draw.Src)
	}

	return dst
}

// fillOffset returns the crop offsets for fill mode based on anchor.
func fillOffset(newW, newH, size int, anchor CropAnchor) (int, int) {
	switch anchor {
	case AnchorTop:
		return (newW - size) / 2, 0
	case AnchorBottom:
		return (newW - size) / 2, newH - size
	case AnchorLeft:
		return 0, (newH - size) / 2
	case AnchorRight:
		return newW - size, (newH - size) / 2
	default: // AnchorCenter
		return (newW - size) / 2, (newH - size) / 2
	}
}

// ResizeTo64 resizes an image to 64×64 using fit mode.
func ResizeTo64(img image.Image) image.Image {
	return Resize(img, PixooSize, ResizeFit)
}
