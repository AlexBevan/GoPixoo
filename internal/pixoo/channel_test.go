package pixoo

import (
	"testing"
)

func TestSetBrightness_Structure(t *testing.T) {
	tests := []struct {
		name  string
		level int
	}{
		{"zero brightness", 0},
		{"medium brightness", 50},
		{"max brightness", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := SetBrightness(tt.level)

			if cmd, ok := payload["Command"].(string); !ok || cmd != CmdSetBrightness {
				t.Errorf("Command = %v, want %s", payload["Command"], CmdSetBrightness)
			}
			if brightness, ok := payload["Brightness"].(int); !ok || brightness != tt.level {
				t.Errorf("Brightness = %v, want %d", payload["Brightness"], tt.level)
			}
			if len(payload) != 2 {
				t.Errorf("payload field count = %d, want 2", len(payload))
			}
		})
	}
}

func TestGetBrightness_Structure(t *testing.T) {
	payload := GetBrightness()

	if cmd, ok := payload["Command"].(string); !ok || cmd != CmdGetDeviceSettings {
		t.Errorf("Command = %v, want %s", payload["Command"], CmdGetDeviceSettings)
	}
	if len(payload) != 1 {
		t.Errorf("payload field count = %d, want 1", len(payload))
	}
}

func TestSetChannel_Structure(t *testing.T) {
	tests := []struct {
		name  string
		index int
	}{
		{"clock channel", 0},
		{"cloud channel", 1},
		{"visualizer channel", 2},
		{"custom channel", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := SetChannel(tt.index)

			if cmd, ok := payload["Command"].(string); !ok || cmd != CmdSetChannel {
				t.Errorf("Command = %v, want %s", payload["Command"], CmdSetChannel)
			}
			if idx, ok := payload["SelectIndex"].(int); !ok || idx != tt.index {
				t.Errorf("SelectIndex = %v, want %d", payload["SelectIndex"], tt.index)
			}
			if len(payload) != 2 {
				t.Errorf("payload field count = %d, want 2", len(payload))
			}
		})
	}
}

func TestSetClockFace_Structure(t *testing.T) {
	tests := []struct {
		name string
		id   int
	}{
		{"first clock face", 0},
		{"mid clock face", 50},
		{"high clock face", 999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := SetClockFace(tt.id)

			if cmd, ok := payload["Command"].(string); !ok || cmd != CmdSetClock {
				t.Errorf("Command = %v, want %s", payload["Command"], CmdSetClock)
			}
			if clkId, ok := payload["ClkId"].(int); !ok || clkId != tt.id {
				t.Errorf("ClkId = %v, want %d", payload["ClkId"], tt.id)
			}
			if len(payload) != 2 {
				t.Errorf("payload field count = %d, want 2", len(payload))
			}
		})
	}
}

func TestGetClockInfo_Structure(t *testing.T) {
	payload := GetClockInfo()

	if cmd, ok := payload["Command"].(string); !ok || cmd != CmdGetClockInfo {
		t.Errorf("Command = %v, want %s", payload["Command"], CmdGetClockInfo)
	}
	if len(payload) != 1 {
		t.Errorf("payload field count = %d, want 1", len(payload))
	}
}
