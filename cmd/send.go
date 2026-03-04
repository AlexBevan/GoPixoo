package cmd

import (
	"fmt"
	"os"

	"github.com/alexbevan/gopixoo/internal/imaging"
	"github.com/alexbevan/gopixoo/internal/pixoo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var sendCmd = &cobra.Command{
	Use:   "send <file>",
	Short: "Send a GIF or image to the Pixoo64 display",
	Long: `Send pushes an image or animated GIF to the Pixoo64 device.

For animated GIFs, all frames are extracted and sent sequentially.
For static images (PNG, JPEG), a single frame is sent.

The image is resized to fit the display using the chosen resize mode.`,
	Args: cobra.ExactArgs(1),
	RunE: runSend,
}

var (
	sendSpeed  int
	sendResize string
	sendSize   int
	sendAnchor string
)

func init() {
	sendCmd.Flags().IntVarP(&sendSpeed, "speed", "s", 100, "frame delay in milliseconds for GIF animation")
	sendCmd.Flags().StringVarP(&sendResize, "resize", "r", "fit", "resize mode: fit, fill (supports --anchor), stretch, or none")
	sendCmd.Flags().IntVar(&sendSize, "size", 64, "target pixel size (default 64 for Pixoo64)")
	sendCmd.Flags().StringVar(&sendAnchor, "anchor", "center", "crop anchor for fill mode: center, top, bottom, left, right")
	rootCmd.AddCommand(sendCmd)
}

func runSend(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	// Validate file exists.
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("file not found: %s", filePath)
	}

	// Validate resize mode.
	mode, err := parseResizeMode(sendResize)
	if err != nil {
		return err
	}

	// Validate anchor.
	anchor, err := parseAnchor(sendAnchor)
	if err != nil {
		return err
	}

	// Validate device IP.
	ip := viper.GetString("device.ip")
	if ip == "" {
		return fmt.Errorf("device IP required: use --ip flag, GOPIXOO_IP env, or config file")
	}

	isVerbose := viper.GetBool("verbose")

	// Extract frames (handles both GIF and static images).
	if isVerbose {
		fmt.Fprintf(os.Stderr, "Loading %s...\n", filePath)
	}
	frames, delays, err := imaging.ExtractGIFFrames(filePath, sendSpeed)
	if err != nil {
		return fmt.Errorf("extract frames: %w", err)
	}

	totalFrames := len(frames)
	if isVerbose {
		fmt.Fprintf(os.Stderr, "Extracted %d frame(s)\n", totalFrames)
	}

	// Resize and encode each frame.
	encodedFrames := make([]string, 0, totalFrames)
	for i, frame := range frames {
		resized := imaging.ResizeWithAnchor(frame, sendSize, mode, anchor)
		encoded := imaging.EncodePixelsSized(resized, sendSize)
		encodedFrames = append(encodedFrames, encoded)
		if isVerbose {
			fmt.Fprintf(os.Stderr, "Encoded frame %d of %d\n", i+1, totalFrames)
		}
	}

	// Connect to device.
	client := pixoo.NewClient(ip)

	// Reset GIF ID on the device.
	if isVerbose {
		fmt.Fprintf(os.Stderr, "Resetting GIF ID on %s...\n", ip)
	}
	resp, err := client.Post(pixoo.ResetGIFID())
	if err != nil {
		return fmt.Errorf("reset gif id: %w", err)
	}
	if isVerbose {
		fmt.Fprintf(os.Stderr, "Reset response: %v\n", resp)
	}

	// Get the current GIF ID to use for this animation.
	resp, err = client.Post(pixoo.GetGIFID())
	if err != nil {
		return fmt.Errorf("get gif id: %w", err)
	}
	picID := 0
	if id, ok := resp["PicId"]; ok {
		if idFloat, ok := id.(float64); ok {
			picID = int(idFloat)
		}
	}
	if isVerbose {
		fmt.Fprintf(os.Stderr, "Using PicID: %d\n", picID)
	}

	// Send each frame sequentially.
	for i, encoded := range encodedFrames {
		speed := sendSpeed
		if i < len(delays) {
			speed = delays[i]
		}

		payload := pixoo.SendGIF(totalFrames, i, picID, speed, sendSize, encoded)
		_, err := client.Post(payload)
		if err != nil {
			return fmt.Errorf("send frame %d: %w", i+1, err)
		}
		if isVerbose {
			fmt.Fprintf(os.Stderr, "Sent frame %d of %d (delay: %dms)\n", i+1, totalFrames, speed)
		}
	}

	fmt.Printf("Sent %d frame(s) to %s\n", totalFrames, ip)
	return nil
}

func parseResizeMode(s string) (imaging.ResizeMode, error) {
	switch s {
	case "fit":
		return imaging.ResizeFit, nil
	case "fill":
		return imaging.ResizeFill, nil
	case "stretch":
		return imaging.ResizeStretch, nil
	case "none":
		return imaging.ResizeNone, nil
	default:
		return "", fmt.Errorf("invalid resize mode %q: must be fit, fill, stretch, or none", s)
	}
}

func parseAnchor(s string) (imaging.CropAnchor, error) {
	switch s {
	case "center":
		return imaging.AnchorCenter, nil
	case "top":
		return imaging.AnchorTop, nil
	case "bottom":
		return imaging.AnchorBottom, nil
	case "left":
		return imaging.AnchorLeft, nil
	case "right":
		return imaging.AnchorRight, nil
	default:
		return "", fmt.Errorf("invalid anchor %q: must be center, top, bottom, left, or right", s)
	}
}
