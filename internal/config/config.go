package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Defaults sets default configuration values.
func Defaults() {
	viper.SetDefault("device.ip", "")
	viper.SetDefault("verbose", false)
}

// DeviceIP returns the configured device IP, or an error if not set.
func DeviceIP() (string, error) {
	ip := viper.GetString("device.ip")
	if ip == "" {
		return "", fmt.Errorf("device IP not set: use --ip flag, GOPIXOO_DEVICE_IP env, or config file")
	}
	return ip, nil
}

// Verbose returns whether verbose output is enabled.
func Verbose() bool {
	return viper.GetBool("verbose")
}

// EnsureConfigDir creates the config directory if it doesn't exist.
func EnsureConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".config", "gopixoo")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}
