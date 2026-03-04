package pixoo

// SetBrightness builds a set-brightness command payload.
func SetBrightness(level int) map[string]interface{} {
	return map[string]interface{}{
		"Command":    CmdSetBrightness,
		"Brightness": level,
	}
}

// GetBrightness builds a get-brightness command payload.
// Uses GetAllConf because the device does not support a dedicated GetBrightness endpoint.
func GetBrightness() map[string]interface{} {
	return map[string]interface{}{
		"Command": CmdGetDeviceSettings,
	}
}

// SetChannel builds a set-channel command payload.
// Index: 0=clock, 1=cloud, 2=visualizer, 3=custom
func SetChannel(index int) map[string]interface{} {
	return map[string]interface{}{
		"Command":    CmdSetChannel,
		"SelectIndex": index,
	}
}

// SetClockFace builds a set-clock-face command payload.
func SetClockFace(id int) map[string]interface{} {
	return map[string]interface{}{
		"Command": CmdSetClock,
		"ClkId":   id,
	}
}

// GetClockInfo builds a get-clock-info command payload.
func GetClockInfo() map[string]interface{} {
	return map[string]interface{}{
		"Command": CmdGetClockInfo,
	}
}
