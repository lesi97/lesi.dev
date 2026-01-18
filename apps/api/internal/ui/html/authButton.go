package html

import (
	"fmt"
	"html"
)

func GenerateAuthButton(url string, messageContent string) string {
	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
			<head>
				<title>%s</title>
			</head>
			<style>
				body {
					background-color: #1c1c1c; 
					width: 100vw; 
					height: 100vh; 
					display: flex; 
					flex-direction: row; 
					align-items: center; 
					justify-content: center; 
					overflow: hidden;
				}
				button {
					border-radius: 4px; 
					padding: 8px;
					cursor: pointer;
				}
			</style>
			<body>
				<a href="%s" target="_blank">
					<button tabindex="-1">%s</button>
				</a>
			</body>
		</html>
		`, messageContent, html.EscapeString(url), messageContent)
	return html
}
