package main

import (
	"fmt"
	"strings"
)

// Format renders a format string using weather conditions.
//
// Format codes (compatible with common wttr.in codes):
//
//	%c  - weather condition icon (emoji)
//	%C  - weather condition description (e.g. "Partly cloudy")
//	%t  - temperature in configured units (e.g. "75°F" or "24°C")
//	%f  - feels-like temperature in configured units
//	%h  - humidity percentage (e.g. "50%")
//	%l  - location name
//	%%  - literal percent sign
//
// Units: "u" for US/imperial (°F), "m" for metric (°C).
func Format(format string, c Conditions, units string) string {
	var result strings.Builder
	i := 0
	for i < len(format) {
		if format[i] != '%' || i+1 >= len(format) {
			result.WriteByte(format[i])
			i++
			continue
		}
		i++ // consume '%'
		switch format[i] {
		case 'c':
			result.WriteString(c.Icon)
		case 'C':
			result.WriteString(c.Description)
		case 't':
			result.WriteString(formatTemp(c.TempF, c.TempC, units))
		case 'f':
			result.WriteString(formatTemp(c.FeelsLikeF, c.FeelsLikeC, units))
		case 'h':
			result.WriteString(fmt.Sprintf("%d%%", c.Humidity))
		case 'l':
			result.WriteString(c.Location)
		case '%':
			result.WriteByte('%')
		default:
			result.WriteByte('%')
			result.WriteByte(format[i])
		}
		i++
	}
	return result.String()
}

func formatTemp(f, c float64, units string) string {
	if units == "m" {
		return fmt.Sprintf("%.0f°C", c)
	}
	return fmt.Sprintf("%.0f°F", f)
}
