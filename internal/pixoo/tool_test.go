package pixoo

import (
	"testing"
)

func TestSetTimer_Structure(t *testing.T) {
	tests := []struct {
		name    string
		minutes int
		seconds int
		status  int
	}{
		{"start 5 minute timer", 5, 0, 1},
		{"start 90 second timer", 1, 30, 1},
		{"stop timer", 0, 0, 0},
		{"zero timer with start", 0, 0, 1},
		{"max timer", 99, 59, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := SetTimer(tt.minutes, tt.seconds, tt.status)

			if cmd, ok := payload["Command"].(string); !ok || cmd != CmdSetCountdown {
				t.Errorf("Command = %v, want %s", payload["Command"], CmdSetCountdown)
			}
			if minute, ok := payload["Minute"].(int); !ok || minute != tt.minutes {
				t.Errorf("Minute = %v, want %d", payload["Minute"], tt.minutes)
			}
			if second, ok := payload["Second"].(int); !ok || second != tt.seconds {
				t.Errorf("Second = %v, want %d", payload["Second"], tt.seconds)
			}
			if status, ok := payload["Status"].(int); !ok || status != tt.status {
				t.Errorf("Status = %v, want %d", payload["Status"], tt.status)
			}
			if len(payload) != 4 {
				t.Errorf("payload field count = %d, want 4", len(payload))
			}
		})
	}
}

func TestSetStopwatch_Structure(t *testing.T) {
	tests := []struct {
		name   string
		status int
		desc   string
	}{
		{"start stopwatch", 1, "start"},
		{"stop stopwatch", 2, "stop"},
		{"reset stopwatch", 0, "reset"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := SetStopwatch(tt.status)

			if cmd, ok := payload["Command"].(string); !ok || cmd != CmdSetStopwatch {
				t.Errorf("Command = %v, want %s", payload["Command"], CmdSetStopwatch)
			}
			if status, ok := payload["Status"].(int); !ok || status != tt.status {
				t.Errorf("Status = %v, want %d", payload["Status"], tt.status)
			}
			if len(payload) != 2 {
				t.Errorf("payload field count = %d, want 2", len(payload))
			}
		})
	}
}

func TestSetScoreboard_Structure(t *testing.T) {
	tests := []struct {
		name string
		blue int
		red  int
	}{
		{"zero score", 0, 0},
		{"blue leading", 5, 3},
		{"red leading", 2, 7},
		{"tie game", 10, 10},
		{"max scores", 999, 999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := SetScoreboard(tt.blue, tt.red)

			if cmd, ok := payload["Command"].(string); !ok || cmd != CmdSetScoreboard {
				t.Errorf("Command = %v, want %s", payload["Command"], CmdSetScoreboard)
			}
			if blue, ok := payload["BlueScore"].(int); !ok || blue != tt.blue {
				t.Errorf("BlueScore = %v, want %d", payload["BlueScore"], tt.blue)
			}
			if red, ok := payload["RedScore"].(int); !ok || red != tt.red {
				t.Errorf("RedScore = %v, want %d", payload["RedScore"], tt.red)
			}
			if len(payload) != 3 {
				t.Errorf("payload field count = %d, want 3", len(payload))
			}
		})
	}
}

func TestSetNoiseMeter_Structure(t *testing.T) {
	tests := []struct {
		name   string
		status int
		desc   string
	}{
		{"start noise meter", 1, "start"},
		{"stop noise meter", 0, "stop"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := SetNoiseMeter(tt.status)

			if cmd, ok := payload["Command"].(string); !ok || cmd != CmdSetNoiseMeter {
				t.Errorf("Command = %v, want %s", payload["Command"], CmdSetNoiseMeter)
			}
			if status, ok := payload["NoiseStatus"].(int); !ok || status != tt.status {
				t.Errorf("NoiseStatus = %v, want %d", payload["NoiseStatus"], tt.status)
			}
			if len(payload) != 2 {
				t.Errorf("payload field count = %d, want 2", len(payload))
			}
		})
	}
}
