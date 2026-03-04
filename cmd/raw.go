package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/alexbevan/gopixoo/internal/pixoo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rawCmd = &cobra.Command{
	Use:     "raw '{\"Command\":\"...\"}'",
	Short:   "Send a raw JSON command to the device",
	Example: `  gopixoo raw '{"Command": "Channel/GetAllConf"}'`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ip := viper.GetString("device.ip")
		if ip == "" {
			return fmt.Errorf("device IP not set (use --ip or config file)")
		}

		var payload map[string]interface{}
		if err := json.Unmarshal([]byte(args[0]), &payload); err != nil {
			return fmt.Errorf("invalid JSON: %w", err)
		}

		client := pixoo.NewClient(ip)
		resp, err := client.Post(payload)
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}

		out, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("format response: %w", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), string(out))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rawCmd)
}
