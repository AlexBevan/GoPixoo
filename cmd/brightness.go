package cmd

import (
	"fmt"
	"strconv"

	"github.com/alexbevan/gopixoo/internal/pixoo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var brightnessCmd = &cobra.Command{
	Use:   "brightness",
	Short: "Get or set display brightness",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var brightnessGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get current display brightness",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ip := viper.GetString("device.ip")
		if ip == "" {
			return fmt.Errorf("device IP not set (use --ip or config file)")
		}

		client := pixoo.NewClient(ip)
		resp, err := client.Post(pixoo.GetBrightness())
		if err != nil {
			return fmt.Errorf("failed to get brightness: %w", err)
		}

		brightness, ok := resp["Brightness"]
		if !ok {
			return fmt.Errorf("unexpected response: missing Brightness field")
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Brightness: %v\n", brightness)
		return nil
	},
}

var brightnessSetCmd = &cobra.Command{
	Use:   "set <0-100>",
	Short: "Set display brightness (0-100)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ip := viper.GetString("device.ip")
		if ip == "" {
			return fmt.Errorf("device IP not set (use --ip or config file)")
		}

		level, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid brightness value: %s (must be 0-100)", args[0])
		}
		if level < 0 || level > 100 {
			return fmt.Errorf("brightness must be between 0 and 100, got %d", level)
		}

		client := pixoo.NewClient(ip)
		if _, err := client.Post(pixoo.SetBrightness(level)); err != nil {
			return fmt.Errorf("failed to set brightness: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Brightness set to %d\n", level)
		return nil
	},
}

func init() {
	brightnessCmd.AddCommand(brightnessGetCmd)
	brightnessCmd.AddCommand(brightnessSetCmd)
	rootCmd.AddCommand(brightnessCmd)
}
