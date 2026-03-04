package cmd

import (
	"fmt"

	"github.com/alexbevan/gopixoo/internal/pixoo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var channelCmd = &cobra.Command{
	Use:   "channel [clock|cloud|visualizer|custom]",
	Short: "Switch display channel",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func newChannelSubCmd(name string, index int) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Switch to %s channel", name),
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ip := viper.GetString("device.ip")
			if ip == "" {
				return fmt.Errorf("device IP not set (use --ip or config file)")
			}

			client := pixoo.NewClient(ip)
			if _, err := client.Post(pixoo.SetChannel(index)); err != nil {
				return fmt.Errorf("failed to switch channel: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Switched to %s channel\n", name)
			return nil
		},
	}
}

func init() {
	channelCmd.AddCommand(newChannelSubCmd("clock", 0))
	channelCmd.AddCommand(newChannelSubCmd("cloud", 1))
	channelCmd.AddCommand(newChannelSubCmd("visualizer", 2))
	channelCmd.AddCommand(newChannelSubCmd("custom", 3))
	rootCmd.AddCommand(channelCmd)
}
