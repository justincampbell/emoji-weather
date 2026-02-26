package main

import "testing"

func TestFormat(t *testing.T) {
	c := Conditions{
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
		{"icon and temp US", "%c %t", "u", "⛅ 75°F"},
		{"icon and temp metric", "%c %t", "m", "⛅ 24°C"},
		{"description", "%C", "u", "Partly cloudy"},
		{"feels like US", "%f", "u", "72°F"},
		{"feels like metric", "%f", "m", "22°C"},
		{"humidity", "%h", "u", "50%"},
		{"location", "%l", "u", "New York"},
		{"literal percent", "100%%", "u", "100%"},
		{"unknown code passthrough", "%x", "u", "%x"},
		{"combined", "%c %t · %h", "u", "⛅ 75°F · 50%"},
		{"trailing percent", "test%", "u", "test%"},
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
