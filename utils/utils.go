package utils

import (
	"fmt"
	"log"
	"strings"
)

type color string

const (
	YELLOW color = "\033[33m"
	GREEN  color = "\033[32m"
	RED    color = "\033[31m"
)

func mapS(c color, s any) string {
	return fmt.Sprintf("%s%s\033[0m", c, s)
}

func mapToColor(c color, args ...any) []string {
	coloredStrings := make([]string, len(args))
	for i, s := range args {
		coloredStrings[i] = mapS(c, s)
	}
	return coloredStrings
}

func LogDebug(args ...any) {
	log.Println(strings.Join(mapToColor(YELLOW, args), " "))
}

func LogGreen(args ...any) {
	log.Println(strings.Join(mapToColor(GREEN, args), " "))
}

func LogRed(args ...any) {
	log.Println(strings.Join(mapToColor(RED, args), " "))
}
