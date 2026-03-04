package pixoo

// SendGIF builds a send-GIF command payload.
// picData is the base64-encoded pixel data for one frame.
func SendGIF(picNum, picOffset, picID, picSpeed, picWidth int, picData string) map[string]interface{} {
	return map[string]interface{}{
		"Command":   CmdSendHTTPGif,
		"PicNum":    picNum,
		"PicWidth":  picWidth,
		"PicOffset": picOffset,
		"PicID":     picID,
		"PicSpeed":  picSpeed,
		"PicData":   picData,
	}
}

// ResetGIFID builds a reset-GIF-ID command payload.
func ResetGIFID() map[string]interface{} {
	return map[string]interface{}{
		"Command": CmdResetHTTPGifID,
	}
}

// GetGIFID builds a get-current-GIF-ID command payload.
func GetGIFID() map[string]interface{} {
	return map[string]interface{}{
		"Command": CmdGetHTTPGifID,
	}
}
