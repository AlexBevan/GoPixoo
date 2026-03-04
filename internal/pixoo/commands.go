package pixoo

// Pixoo64 API command names.
const (
	CmdGetDeviceSettings   = "Channel/GetAllConf"
	CmdSetBrightness       = "Channel/SetBrightness"
	CmdGetBrightness       = "Channel/GetBrightness"
	CmdSetScreen           = "Channel/OnOffScreen"
	CmdSetChannel          = "Channel/SetIndex"
	CmdSetClock            = "Channel/SetClockSelectId"
	CmdGetClockInfo        = "Channel/GetClockInfo"
	CmdSetCustomPage       = "Channel/SetCustomPageIndex"
	CmdSetVisualizer       = "Channel/SetEqPosition"
	CmdSendHTTPText        = "Draw/SendHttpText"
	CmdSendHTTPGif         = "Draw/SendHttpGif"
	CmdResetHTTPGifID      = "Draw/ResetHttpGifId"
	CmdClearHTTPText       = "Draw/ClearHttpText"
	CmdSendHTTPItemList    = "Draw/SendHttpItemList"
	CmdGetHTTPGifID        = "Draw/GetHttpGifId"
	CmdCommandList         = "Draw/CommandList"
	CmdSetCountdown        = "Tools/SetTimer"
	CmdSetStopwatch        = "Tools/SetStopWatch"
	CmdSetScoreboard       = "Tools/SetScoreBoard"
	CmdSetNoiseMeter       = "Tools/SetNoiseStatus"
	CmdSetBuzzer           = "Device/PlayBuzzer"
	CmdReboot              = "Device/SysReboot"
	CmdGetDeviceTime       = "Device/GetDeviceTime"
)
