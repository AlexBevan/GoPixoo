package cmd

import (
	"fmt"
	"strconv"

	"github.com/alexbevan/gopixoo/internal/pixoo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const drawSize = 64

var drawCmd = &cobra.Command{
	Use:   "draw [pixel|fill|clear]",
	Short: "Draw primitives on the display",
	Long:  `Draw pixels and shapes directly on the Pixoo64 display.`,
}

var drawPixelCmd = &cobra.Command{
	Use:   "pixel <x> <y>",
	Short: "Set a single pixel on the display",
	Args:  cobra.ExactArgs(2),
	RunE:  runDrawPixel,
}

var drawFillCmd = &cobra.Command{
	Use:   "fill",
	Short: "Fill the entire display with a solid color",
	RunE:  runDrawFill,
}

var drawClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear the display (fill with black)",
	RunE:  runDrawClear,
}

var drawColor string

func init() {
	drawPixelCmd.Flags().StringVar(&drawColor, "color", "#FFFFFF", "pixel color in hex (e.g. #FF0000)")
	drawFillCmd.Flags().StringVar(&drawColor, "color", "#FF0000", "fill color in hex (e.g. #00FF00)")

	drawCmd.AddCommand(drawPixelCmd)
	drawCmd.AddCommand(drawFillCmd)
	drawCmd.AddCommand(drawClearCmd)
	rootCmd.AddCommand(drawCmd)
}

// sendFrame resets the GIF ID and sends a single frame to the device.
func sendFrame(ip string, picData string) error {
	client := pixoo.NewClient(ip)

	if _, err := client.Post(pixoo.ResetGIFID()); err != nil {
		return fmt.Errorf("reset gif id: %w", err)
	}

	resp, err := client.Post(pixoo.GetGIFID())
	if err != nil {
		return fmt.Errorf("get gif id: %w", err)
	}
	picID := 0
	if id, ok := resp["PicId"]; ok {
		if f, ok := id.(float64); ok {
			picID = int(f)
		}
	}

	payload := pixoo.SendGIF(1, 0, picID, 1000, drawSize, picData)
	if _, err := client.Post(payload); err != nil {
		return fmt.Errorf("send frame: %w", err)
	}
	return nil
}

func requireIP() (string, error) {
	ip := viper.GetString("device.ip")
	if ip == "" {
		return "", fmt.Errorf("device IP required: use --ip flag, GOPIXOO_IP env, or config file")
	}
	return ip, nil
}

func runDrawPixel(cmd *cobra.Command, args []string) error {
	x, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid x coordinate: %w", err)
	}
	y, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid y coordinate: %w", err)
	}
	if x < 0 || x >= drawSize || y < 0 || y >= drawSize {
		return fmt.Errorf("coordinates out of range: must be 0-%d", drawSize-1)
	}

	color, _ := cmd.Flags().GetString("color")
	r, g, b, err := pixoo.ParseHexColor(color)
	if err != nil {
		return err
	}

	ip, err := requireIP()
	if err != nil {
		return err
	}

	picData := pixoo.BuildPixelFrame(drawSize, x, y, r, g, b)
	if err := sendFrame(ip, picData); err != nil {
		return err
	}
	fmt.Printf("Drew pixel at (%d, %d) color %s on %s\n", x, y, color, ip)
	return nil
}

func runDrawFill(cmd *cobra.Command, args []string) error {
	color, _ := cmd.Flags().GetString("color")
	r, g, b, err := pixoo.ParseHexColor(color)
	if err != nil {
		return err
	}

	ip, err := requireIP()
	if err != nil {
		return err
	}

	picData := pixoo.BuildSolidFrame(drawSize, r, g, b)
	if err := sendFrame(ip, picData); err != nil {
		return err
	}
	fmt.Printf("Filled display with %s on %s\n", color, ip)
	return nil
}

func runDrawClear(cmd *cobra.Command, args []string) error {
	ip, err := requireIP()
	if err != nil {
		return err
	}

	picData := pixoo.BuildSolidFrame(drawSize, 0, 0, 0)
	if err := sendFrame(ip, picData); err != nil {
		return err
	}
	fmt.Printf("Cleared display on %s\n", ip)
	return nil
}
