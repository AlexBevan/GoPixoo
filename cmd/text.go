package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alexbevan/gopixoo/internal/pixoo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var textCmd = &cobra.Command{
	Use:   "text",
	Short: "Send text to display",
	Long:  `Display scrolling text on the Pixoo64. Use "text send" to show text or "text clear" to remove it.`,
}

var textSendCmd = &cobra.Command{
	Use:   "send <message>",
	Short: "Send text to the Pixoo64 display",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runTextSend,
}

var textClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all text from the display",
	RunE:  runTextClear,
}

func init() {
	textSendCmd.Flags().Int("x", 0, "X position")
	textSendCmd.Flags().Int("y", 0, "Y position")
	textSendCmd.Flags().Int("font", 0, "font index (0-7)")
	textSendCmd.Flags().String("color", "#FFFFFF", "text color as hex (e.g. #FF0000)")
	textSendCmd.Flags().Int("speed", 100, "scroll speed in ms")
	textSendCmd.Flags().Int("dir", 0, "scroll direction (0=left, 1=right)")
	textSendCmd.Flags().Int("align", 0, "static alignment (1=left, 2=center, 3=right); 0 omits the field to allow scrolling")
	textSendCmd.Flags().Int("id", 1, "text ID for managing multiple texts")
	textSendCmd.Flags().Int("width", 64, "text width in pixels")

	textCmd.AddCommand(textSendCmd)
	textCmd.AddCommand(textClearCmd)
	rootCmd.AddCommand(textCmd)
}

func runTextSend(cmd *cobra.Command, args []string) error {
	ip := viper.GetString("device.ip")
	if ip == "" {
		return fmt.Errorf("device IP required: use --ip flag, GOPIXOO_IP env, or config file")
	}

	message := strings.Join(args, " ")

	x, _ := cmd.Flags().GetInt("x")
	y, _ := cmd.Flags().GetInt("y")
	font, _ := cmd.Flags().GetInt("font")
	color, _ := cmd.Flags().GetString("color")
	speed, _ := cmd.Flags().GetInt("speed")
	dir, _ := cmd.Flags().GetInt("dir")
	align, _ := cmd.Flags().GetInt("align")
	id, _ := cmd.Flags().GetInt("id")
	width, _ := cmd.Flags().GetInt("width")

	client := pixoo.NewClient(ip)
	payload := pixoo.SendText(id, x, y, dir, font, width, message, color, speed, align)

	isVerbose := viper.GetBool("verbose")
	if isVerbose {
		fmt.Fprintf(os.Stderr, "Sending text %q to %s...\n", message, ip)
	}

	resp, err := client.Post(payload)
	if err != nil {
		return fmt.Errorf("send text: %w", err)
	}

	if isVerbose {
		fmt.Fprintf(os.Stderr, "Response: %v\n", resp)
	}

	fmt.Printf("Text sent to %s: %q\n", ip, message)
	return nil
}

func runTextClear(cmd *cobra.Command, args []string) error {
	ip := viper.GetString("device.ip")
	if ip == "" {
		return fmt.Errorf("device IP required: use --ip flag, GOPIXOO_IP env, or config file")
	}

	client := pixoo.NewClient(ip)
	payload := pixoo.ClearText()

	isVerbose := viper.GetBool("verbose")
	if isVerbose {
		fmt.Fprintf(os.Stderr, "Clearing text on %s...\n", ip)
	}

	resp, err := client.Post(payload)
	if err != nil {
		return fmt.Errorf("clear text: %w", err)
	}

	if isVerbose {
		fmt.Fprintf(os.Stderr, "Response: %v\n", resp)
	}

	fmt.Printf("Text cleared on %s\n", ip)
	return nil
}
