package pixoo

// SendText builds a payload to display scrolling text on the Pixoo64.
// The align parameter is only included when non-zero; the Pixoo firmware
// treats its presence as a request for static positioning which prevents
// the text from scrolling.
func SendText(id, x, y, dir, font, width int, text, color string, speed, align int) map[string]interface{} {
	payload := map[string]interface{}{
		"Command":    CmdSendHTTPText,
		"TextId":     id,
		"x":          x,
		"y":          y,
		"dir":        dir,
		"font":       font,
		"TextWidth":  width,
		"TextString": text,
		"speed":      speed,
		"color":      color,
	}
	if align != 0 {
		payload["align"] = align
	}
	return payload
}

// ClearText builds a payload to clear all HTTP text from the display.
func ClearText() map[string]interface{} {
	return map[string]interface{}{
		"Command": CmdClearHTTPText,
	}
}
