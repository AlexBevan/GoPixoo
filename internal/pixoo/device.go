package pixoo

// SetScreenOn builds a screen on/off command payload.
func SetScreenOn(on bool) map[string]interface{} {
	val := 0
	if on {
		val = 1
	}
	return map[string]interface{}{
		"Command": CmdSetScreen,
		"OnOff":   val,
	}
}

// GetDeviceSettings builds a get-all-settings command payload.
func GetDeviceSettings() map[string]interface{} {
	return map[string]interface{}{
		"Command": CmdGetDeviceSettings,
	}
}

// Reboot builds a reboot command payload.
func Reboot() map[string]interface{} {
	return map[string]interface{}{
		"Command": CmdReboot,
	}
}

// GetDeviceTime builds a get-device-time command payload.
func GetDeviceTime() map[string]interface{} {
	return map[string]interface{}{
		"Command": CmdGetDeviceTime,
	}
}
