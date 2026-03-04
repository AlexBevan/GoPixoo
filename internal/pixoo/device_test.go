package pixoo

import (
	"testing"
)

func TestSetScreenOn_Structure(t *testing.T) {
	tests := []struct {
		name     string
		on       bool
		wantOnOff int
	}{
		{"screen on", true, 1},
		{"screen off", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := SetScreenOn(tt.on)

			if cmd, ok := payload["Command"].(string); !ok || cmd != CmdSetScreen {
				t.Errorf("Command = %v, want %s", payload["Command"], CmdSetScreen)
			}
			if onOff, ok := payload["OnOff"].(int); !ok || onOff != tt.wantOnOff {
				t.Errorf("OnOff = %v, want %d", payload["OnOff"], tt.wantOnOff)
			}
			if len(payload) != 2 {
				t.Errorf("payload field count = %d, want 2", len(payload))
			}
		})
	}
}

func TestGetDeviceSettings_Structure(t *testing.T) {
	payload := GetDeviceSettings()

	if cmd, ok := payload["Command"].(string); !ok || cmd != CmdGetDeviceSettings {
		t.Errorf("Command = %v, want %s", payload["Command"], CmdGetDeviceSettings)
	}
	if len(payload) != 1 {
		t.Errorf("payload field count = %d, want 1", len(payload))
	}
}

func TestReboot_Structure(t *testing.T) {
	payload := Reboot()

	if cmd, ok := payload["Command"].(string); !ok || cmd != CmdReboot {
		t.Errorf("Command = %v, want %s", payload["Command"], CmdReboot)
	}
	if len(payload) != 1 {
		t.Errorf("payload field count = %d, want 1", len(payload))
	}
}

func TestGetDeviceTime_Structure(t *testing.T) {
	payload := GetDeviceTime()

	if cmd, ok := payload["Command"].(string); !ok || cmd != CmdGetDeviceTime {
		t.Errorf("Command = %v, want %s", payload["Command"], CmdGetDeviceTime)
	}
	if len(payload) != 1 {
		t.Errorf("payload field count = %d, want 1", len(payload))
	}
}
