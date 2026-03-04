package cmd

import (
	"fmt"

	"github.com/alexbevan/gopixoo/internal/pixoo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var displayCmd = &cobra.Command{
	Use:   "display",
	Short: "Turn the display on or off",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var displayOnCmd = &cobra.Command{
	Use:   "on",
	Short: "Turn the display on",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ip := viper.GetString("device.ip")
		if ip == "" {
			return fmt.Errorf("device IP not set (use --ip or config file)")
		}

		client := pixoo.NewClient(ip)
		if _, err := client.Post(pixoo.SetScreenOn(true)); err != nil {
			return fmt.Errorf("failed to turn display on: %w", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Display turned on")
		return nil
	},
}

var displayOffCmd = &cobra.Command{
	Use:   "off",
	Short: "Turn the display off",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ip := viper.GetString("device.ip")
		if ip == "" {
			return fmt.Errorf("device IP not set (use --ip or config file)")
		}

		client := pixoo.NewClient(ip)
		if _, err := client.Post(pixoo.SetScreenOn(false)); err != nil {
			return fmt.Errorf("failed to turn display off: %w", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Display turned off")
		return nil
	},
}

func init() {
	displayCmd.AddCommand(displayOnCmd)
	displayCmd.AddCommand(displayOffCmd)
	rootCmd.AddCommand(displayCmd)
}
