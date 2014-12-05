package main

import "fmt"

var conditionIcons = map[string]string{
	"clear-day": "☀️",
}

func main() {
	fmt.Println(formatConditions("clear-day"))
}

func formatConditions(condition string) (icon string) {
	icon, ok := conditionIcons[condition]
	if !ok {
		icon = condition
	}
	return
}
