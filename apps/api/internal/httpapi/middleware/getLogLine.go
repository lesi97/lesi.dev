package middleware

import (
	"fmt"
	"time"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func GetLogLine(path string, status int, duration time.Duration, nonceElapsed time.Duration, nonceOk bool) string {
	pathColour := utils.Colours["brightBlue"]
	timeColour := utils.Colours["green"]
	statusColour := utils.Colours["green"]

	if duration > 100*time.Millisecond {
		timeColour = utils.Colours["brightRed"] + utils.Colours["bold"]
	}

	if status >= 400 {
		statusColour = utils.Colours["brightRed"] + utils.Colours["bold"]
	}

	reset := utils.Colours["reset"]

	statusBlock := fmt.Sprintf("%vStatus: %v%v", statusColour, status, reset)
	pathBlock := fmt.Sprintf("%v%v%v", pathColour, path, reset)
	timeBlock := fmt.Sprintf("%vtook %v%v", timeColour, duration, reset)
	nonceBlock := ""
	if nonceOk {
		nonceBlock = fmt.Sprintf("%vnonce %v%v", timeColour, FormatNonceElapsed(nonceElapsed), reset)
	}
	logLine := fmt.Sprintf("%v | %v %v", statusBlock, pathBlock, timeBlock)
	if nonceBlock != "" {
		logLine = fmt.Sprintf("%v | %v", logLine, nonceBlock)
	}

	return logLine
}
