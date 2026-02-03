package middleware

import (
	"fmt"
	"net/http"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func LogNotFoundRequest(
	logger *utils.Logger,
	logLine string,
	header http.Header,
	hasNightbotHeaders bool,
	hasStreamElementsHeader bool,
) {
	if hasNightbotHeaders || hasStreamElementsHeader {
		return
	}

	fmt.Println()
	logger.Printf("%v", logLine)
	reset := utils.Colours["reset"]
	for key, values := range header {
		for _, value := range values {
			keyBlock := fmt.Sprintf("%v%v%v: %v", utils.Colours["brightBlue"], utils.Colours["dim"], key, reset)
			logger.Printf("%v%v\n", keyBlock, value)
		}
	}
	fmt.Println()
}
