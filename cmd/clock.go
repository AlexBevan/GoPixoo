package cmd

import (
	"fmt"
	"strconv"

	"github.com/alexbevan/gopixoo/internal/pixoo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var clockCmd = &cobra.Command{
	Use:   "clock",
	Short: "Get or set the clock face",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var clockSetCmd = &cobra.Command{
	Use:   "set <id>",
	Short: "Set the clock face by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ip := viper.GetString("device.ip")
		if ip == "" {
			return fmt.Errorf("device IP not set (use --ip or config file)")
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid clock face ID: %s (must be an integer)", args[0])
		}

		client := pixoo.NewClient(ip)
		if _, err := client.Post(pixoo.SetClockFace(id)); err != nil {
			return fmt.Errorf("failed to set clock face: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Clock face set to %d\n", id)
		return nil
	},
}

var clockGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get current clock face info",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ip := viper.GetString("device.ip")
		if ip == "" {
			return fmt.Errorf("device IP not set (use --ip or config file)")
		}

		client := pixoo.NewClient(ip)
		resp, err := client.Post(pixoo.GetClockInfo())
		if err != nil {
			return fmt.Errorf("failed to get clock info: %w", err)
		}

		clkID, ok := resp["ClockId"]
		if !ok {
			// Try alternate key
			clkID, ok = resp["ClkId"]
		}
		if ok {
			fmt.Fprintf(cmd.OutOrStdout(), "Current clock face ID: %v\n", clkID)
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "Clock info: %v\n", resp)
		}
		return nil
	},
}

func init() {
	clockCmd.AddCommand(clockSetCmd)
	clockCmd.AddCommand(clockGetCmd)
	rootCmd.AddCommand(clockCmd)
}
