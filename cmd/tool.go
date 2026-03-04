package cmd

import (
	"fmt"
	"strconv"

	"github.com/alexbevan/gopixoo/internal/pixoo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var toolCmd = &cobra.Command{
	Use:   "tool",
	Short: "Pixoo built-in tools",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var timerStop bool

var timerCmd = &cobra.Command{
	Use:   "timer <minutes> <seconds>",
	Short: "Start or stop a countdown timer",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ip := viper.GetString("device.ip")
		if ip == "" {
			return fmt.Errorf("device IP not set (use --ip or config file)")
		}

		minutes, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid minutes: %s", args[0])
		}
		seconds, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid seconds: %s", args[1])
		}

		status := 1
		if timerStop {
			status = 0
		}

		client := pixoo.NewClient(ip)
		if _, err := client.Post(pixoo.SetTimer(minutes, seconds, status)); err != nil {
			return fmt.Errorf("failed to set timer: %w", err)
		}

		if timerStop {
			fmt.Fprintln(cmd.OutOrStdout(), "Timer stopped")
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "Timer started: %dm %ds\n", minutes, seconds)
		}
		return nil
	},
}

var stopwatchCmd = &cobra.Command{
	Use:   "stopwatch <start|stop|reset>",
	Short: "Control the stopwatch",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ip := viper.GetString("device.ip")
		if ip == "" {
			return fmt.Errorf("device IP not set (use --ip or config file)")
		}

		var status int
		switch args[0] {
		case "start":
			status = 1
		case "stop":
			status = 2
		case "reset":
			status = 0
		default:
			return fmt.Errorf("invalid action %q: must be start, stop, or reset", args[0])
		}

		client := pixoo.NewClient(ip)
		if _, err := client.Post(pixoo.SetStopwatch(status)); err != nil {
			return fmt.Errorf("failed to set stopwatch: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Stopwatch %s\n", args[0])
		return nil
	},
}

var scoreboardCmd = &cobra.Command{
	Use:   "scoreboard <blue> <red>",
	Short: "Set scoreboard scores",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ip := viper.GetString("device.ip")
		if ip == "" {
			return fmt.Errorf("device IP not set (use --ip or config file)")
		}

		blue, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid blue score: %s", args[0])
		}
		red, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid red score: %s", args[1])
		}

		client := pixoo.NewClient(ip)
		if _, err := client.Post(pixoo.SetScoreboard(blue, red)); err != nil {
			return fmt.Errorf("failed to set scoreboard: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Scoreboard: blue %d — red %d\n", blue, red)
		return nil
	},
}

var noiseCmd = &cobra.Command{
	Use:   "noise <start|stop>",
	Short: "Toggle noise meter",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ip := viper.GetString("device.ip")
		if ip == "" {
			return fmt.Errorf("device IP not set (use --ip or config file)")
		}

		var status int
		switch args[0] {
		case "start":
			status = 1
		case "stop":
			status = 0
		default:
			return fmt.Errorf("invalid action %q: must be start or stop", args[0])
		}

		client := pixoo.NewClient(ip)
		if _, err := client.Post(pixoo.SetNoiseMeter(status)); err != nil {
			return fmt.Errorf("failed to set noise meter: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Noise meter %s\n", args[0])
		return nil
	},
}

func init() {
	timerCmd.Flags().BoolVar(&timerStop, "stop", false, "stop the timer")
	toolCmd.AddCommand(timerCmd)
	toolCmd.AddCommand(stopwatchCmd)
	toolCmd.AddCommand(scoreboardCmd)
	toolCmd.AddCommand(noiseCmd)
	rootCmd.AddCommand(toolCmd)
}
