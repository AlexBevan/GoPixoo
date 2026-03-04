package cmd

import (
	"fmt"
	"os"

	"github.com/alexbevan/gopixoo/internal/imaging"
	"github.com/alexbevan/gopixoo/internal/pixoo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	imageResize string
	imageSize   int
	imageAnchor string
)

var imageCmd = &cobra.Command{
	Use:   "image <file>",
	Short: "Push a static image to the Pixoo64 display",
	Long: `Push a single static image (PNG, JPEG, BMP, or GIF first frame) to the Pixoo64.

The image is resized to the target pixel size using the chosen resize mode,
then sent as a single frame via the device GIF protocol.`,
	Args: cobra.ExactArgs(1),
	RunE: runImage,
}

func init() {
	imageCmd.Flags().StringVar(&imageResize, "resize", "fit", "resize mode: fit, fill (supports --anchor), stretch, or none")
	imageCmd.Flags().IntVar(&imageSize, "size", 64, "target pixel size (default 64 for Pixoo64)")
	imageCmd.Flags().StringVar(&imageAnchor, "anchor", "center", "crop anchor for fill mode: center, top, bottom, left, right")
	rootCmd.AddCommand(imageCmd)
}

func runImage(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("file not found: %s", filePath)
	}

	mode, err := parseResizeMode(imageResize)
	if err != nil {
		return err
	}

	anchor, err := parseAnchor(imageAnchor)
	if err != nil {
		return err
	}

	ip := viper.GetString("device.ip")
	if ip == "" {
		return fmt.Errorf("device IP required: use --ip flag, GOPIXOO_IP env, or config file")
	}

	isVerbose := viper.GetBool("verbose")

	// Load the image (GIF decoded as first frame by image.Decode).
	if isVerbose {
		fmt.Fprintf(os.Stderr, "Loading %s...\n", filePath)
	}
	img, err := imaging.LoadImage(filePath)
	if err != nil {
		return fmt.Errorf("load image: %w", err)
	}

	// Resize and encode pixels.
	resized := imaging.ResizeWithAnchor(img, imageSize, mode, anchor)
	encoded := imaging.EncodePixelsSized(resized, imageSize)
	if isVerbose {
		fmt.Fprintf(os.Stderr, "Image encoded (%s, size=%d)\n", imageResize, imageSize)
	}

	// Send single frame to device.
	client := pixoo.NewClient(ip)

	if isVerbose {
		fmt.Fprintf(os.Stderr, "Resetting GIF ID on %s...\n", ip)
	}
	if _, err := client.Post(pixoo.ResetGIFID()); err != nil {
		return fmt.Errorf("reset gif id: %w", err)
	}

	resp, err := client.Post(pixoo.GetGIFID())
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

	payload := pixoo.SendGIF(1, 0, picID, 0, imageSize, encoded)
	if _, err := client.Post(payload); err != nil {
		return fmt.Errorf("send image: %w", err)
	}

	fmt.Printf("Image sent to %s (%dx%d, %s)\n", ip, imageSize, imageSize, imageResize)
	return nil
}
