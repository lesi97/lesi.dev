package utils

import (
	"context"
	"time"
)


func (l *Logger) LogExecutionTime(name string, start time.Time, ctx context.Context) {
	pathColour := Colours["brightBlack"]
	timeColour := Colours["green"]
	duration := time.Since(start)

	if ctx != nil {
		select {
		case <-ctx.Done():
			return
		default:
		}
	}

	if duration > 100*time.Millisecond {
		timeColour = Colours["brightRed"] + Colours["bold"]
	}

	l.Printf("%s%s %stook %s%s\n", pathColour, name, timeColour, duration, Colours["reset"])
}
