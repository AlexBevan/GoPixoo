package pixoo

// SetTimer builds a countdown timer payload.
// Status: 1=start, 0=stop.
func SetTimer(minutes, seconds, status int) map[string]interface{} {
	return map[string]interface{}{
		"Command": CmdSetCountdown,
		"Minute":  minutes,
		"Second":  seconds,
		"Status":  status,
	}
}

// SetStopwatch builds a stopwatch payload.
// Status: 1=start, 2=stop, 0=reset.
func SetStopwatch(status int) map[string]interface{} {
	return map[string]interface{}{
		"Command": CmdSetStopwatch,
		"Status":  status,
	}
}

// SetScoreboard builds a scoreboard payload.
func SetScoreboard(blue, red int) map[string]interface{} {
	return map[string]interface{}{
		"Command":   CmdSetScoreboard,
		"BlueScore": blue,
		"RedScore":  red,
	}
}

// SetNoiseMeter builds a noise meter payload.
// Status: 1=start, 0=stop.
func SetNoiseMeter(status int) map[string]interface{} {
	return map[string]interface{}{
		"Command":     CmdSetNoiseMeter,
		"NoiseStatus": status,
	}
}
