package main

import (
	"testing"

	"github.com/justincampbell/emoji-weather/providers"
)

func TestFormat(t *testing.T) {
	c := providers.Conditions{
		Icon:        "⛅",
		Description: "Partly cloudy",
		TempF:       75.0,
		TempC:       24.0,
		FeelsLikeF:  72.0,
		FeelsLikeC:  22.0,
		Humidity:    50,
		Location:    "New York",
	}

	tests := []struct {
		name   string
		format string
		units  string
		want   string
	}{
		{"icon and temp US", "%c %t", "f", "⛅ 75°F"},
		{"icon and temp metric", "%c %t", "c", "⛅ 24°C"},
		{"description", "%C", "f", "Partly cloudy"},
		{"feels like US", "%f", "f", "72°F"},
		{"feels like metric", "%f", "c", "22°C"},
		{"humidity", "%h", "f", "50%"},
		{"location", "%l", "f", "New York"},
		{"literal percent", "100%%", "f", "100%"},
		{"unknown code passthrough", "%x", "f", "%x"},
		{"combined", "%c %t · %h", "f", "⛅ 75°F · 50%"},
		{"trailing percent", "test%", "f", "test%"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Format(tt.format, c, tt.units)
			if got != tt.want {
				t.Errorf("Format(%q, ..., %q) = %q, want %q", tt.format, tt.units, got, tt.want)
			}
		})
	}
}
