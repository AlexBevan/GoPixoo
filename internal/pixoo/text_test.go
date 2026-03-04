package pixoo

import (
	"testing"
)

func TestSendText_Structure(t *testing.T) {
	tests := []struct {
		name  string
		id    int
		x     int
		y     int
		dir   int
		font  int
		width int
		text  string
		color string
		speed int
		align int
	}{
		{
			name:  "basic text",
			id:    1,
			x:     10,
			y:     20,
			dir:   0,
			font:  2,
			width: 64,
			text:  "Hello",
			color: "#FF0000",
			speed: 50,
			align: 1,
		},
		{
			name:  "zero values",
			id:    0,
			x:     0,
			y:     0,
			dir:   0,
			font:  0,
			width: 0,
			text:  "",
			color: "",
			speed: 0,
			align: 0,
		},
		{
			name:  "max values",
			id:    999,
			x:     64,
			y:     64,
			dir:   3,
			font:  7,
			width: 128,
			text:  "Long scrolling text message",
			color: "#00FFFF",
			speed: 100,
			align: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := SendText(tt.id, tt.x, tt.y, tt.dir, tt.font, tt.width, tt.text, tt.color, tt.speed, tt.align)

			if cmd, ok := payload["Command"].(string); !ok || cmd != CmdSendHTTPText {
				t.Errorf("Command = %v, want %s", payload["Command"], CmdSendHTTPText)
			}
			if textId, ok := payload["TextId"].(int); !ok || textId != tt.id {
				t.Errorf("TextId = %v, want %d", payload["TextId"], tt.id)
			}
			if xPos, ok := payload["x"].(int); !ok || xPos != tt.x {
				t.Errorf("x = %v, want %d", payload["x"], tt.x)
			}
			if yPos, ok := payload["y"].(int); !ok || yPos != tt.y {
				t.Errorf("y = %v, want %d", payload["y"], tt.y)
			}
			if dir, ok := payload["dir"].(int); !ok || dir != tt.dir {
				t.Errorf("dir = %v, want %d", payload["dir"], tt.dir)
			}
			if font, ok := payload["font"].(int); !ok || font != tt.font {
				t.Errorf("font = %v, want %d", payload["font"], tt.font)
			}
			if width, ok := payload["TextWidth"].(int); !ok || width != tt.width {
				t.Errorf("TextWidth = %v, want %d", payload["TextWidth"], tt.width)
			}
			if textStr, ok := payload["TextString"].(string); !ok || textStr != tt.text {
				t.Errorf("TextString = %v, want %q", payload["TextString"], tt.text)
			}
			if color, ok := payload["color"].(string); !ok || color != tt.color {
				t.Errorf("color = %v, want %q", payload["color"], tt.color)
			}
			if speed, ok := payload["speed"].(int); !ok || speed != tt.speed {
				t.Errorf("speed = %v, want %d", payload["speed"], tt.speed)
			}
			if align, ok := payload["align"].(int); !ok || align != tt.align {
				t.Errorf("align = %v, want %d", payload["align"], tt.align)
			}
			if len(payload) != 11 {
				t.Errorf("payload field count = %d, want 11", len(payload))
			}
		})
	}
}

func TestClearText_Structure(t *testing.T) {
	payload := ClearText()

	if cmd, ok := payload["Command"].(string); !ok || cmd != CmdClearHTTPText {
		t.Errorf("Command = %v, want %s", payload["Command"], CmdClearHTTPText)
	}
	if len(payload) != 1 {
		t.Errorf("payload field count = %d, want 1", len(payload))
	}
}
