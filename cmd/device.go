package cmd

import (
	"fmt"
	"io"

	"github.com/alexbevan/gopixoo/internal/pixoo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Device information and management",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var deviceInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show device settings and configuration",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ip := viper.GetString("device.ip")
		if ip == "" {
			return fmt.Errorf("device IP not set (use --ip or config file)")
		}

		client := pixoo.NewClient(ip)
		resp, err := client.Post(pixoo.GetDeviceSettings())
		if err != nil {
			return fmt.Errorf("failed to get device settings: %w", err)
		}

		w := cmd.OutOrStdout()
		fmt.Fprintf(w, "Device Settings (%s)\n", ip)
		fmt.Fprintln(w, "─────────────────────────────")
		printField(w, "Brightness", resp, "Brightness")
		printField(w, "Light Switch", resp, "LightSwitch")
		printField(w, "Clock Face ID", resp, "CurClockId")
		printField(w, "Clock Rotation", resp, "ClockTime")
		printField(w, "Power-On Channel", resp, "PowerOnChannelId")
		printField(w, "Gallery Time", resp, "GalleryTime")
		printField(w, "Temperature Mode", resp, "TemperatureMode")
		printField(w, "Time 24h Format", resp, "Time24Flag")
		printField(w, "Gyro Angle", resp, "GyrateAngle")
		printField(w, "Rotation Flag", resp, "RotationFlag")
		printField(w, "Mirror Mode", resp, "MirrorFlag")
		return nil
	},
}

var deviceTimeCmd = &cobra.Command{
	Use:   "time",
	Short: "Show device time",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ip := viper.GetString("device.ip")
		if ip == "" {
			return fmt.Errorf("device IP not set (use --ip or config file)")
		}

		client := pixoo.NewClient(ip)
		resp, err := client.Post(pixoo.GetDeviceTime())
		if err != nil {
			return fmt.Errorf("failed to get device time: %w", err)
		}

		w := cmd.OutOrStdout()
		fmt.Fprintf(w, "Device Time (%s)\n", ip)
		fmt.Fprintln(w, "─────────────────────────────")
		if utc, ok := resp["UTCTime"]; ok {
			fmt.Fprintf(w, "  %-18s %.0f\n", "UTC Time:", utc)
		}
		printField(w, "Local Time", resp, "LocalTime")
		return nil
	},
}

var deviceRebootCmd = &cobra.Command{
	Use:   "reboot",
	Short: "Reboot the device",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ip := viper.GetString("device.ip")
		if ip == "" {
			return fmt.Errorf("device IP not set (use --ip or config file)")
		}

		client := pixoo.NewClient(ip)
		if _, err := client.Post(pixoo.Reboot()); err != nil {
			return fmt.Errorf("failed to reboot device: %w", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Device rebooting...")
		return nil
	},
}

func printField(w io.Writer, label string, resp map[string]interface{}, key string) {
	if val, ok := resp[key]; ok {
		fmt.Fprintf(w, "  %-18s %v\n", label+":", val)
	}
}

func init() {
	deviceCmd.AddCommand(deviceInfoCmd)
	deviceCmd.AddCommand(deviceTimeCmd)
	deviceCmd.AddCommand(deviceRebootCmd)
	rootCmd.AddCommand(deviceCmd)
}
