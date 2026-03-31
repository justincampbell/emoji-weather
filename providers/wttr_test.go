package providers

import "testing"

const sampleWttrJSON = `{
  "current_condition": [{
    "FeelsLikeC": "22",
    "FeelsLikeF": "72",
    "humidity": "50",
    "temp_C": "24",
    "temp_F": "75",
    "weatherCode": "116",
    "weatherDesc": [{"value": "Partly cloudy"}]
  }],
  "nearest_area": [{
    "areaName": [{"value": "New York"}],
    "country": [{"value": "United States of America"}]
  }]
}`

func TestParseWttrJSON(t *testing.T) {
	c, err := parseWttrJSON([]byte(sampleWttrJSON))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if c.Icon != "⛅" {
		t.Errorf("Icon = %q, want ⛅", c.Icon)
	}
	if c.Description != "Partly cloudy" {
		t.Errorf("Description = %q, want 'Partly cloudy'", c.Description)
	}
	if c.TempF != 75 {
		t.Errorf("TempF = %v, want 75", c.TempF)
	}
	if c.TempC != 24 {
		t.Errorf("TempC = %v, want 24", c.TempC)
	}
	if c.FeelsLikeF != 72 {
		t.Errorf("FeelsLikeF = %v, want 72", c.FeelsLikeF)
	}
	if c.Humidity != 50 {
		t.Errorf("Humidity = %v, want 50", c.Humidity)
	}
	if c.Location != "New York" {
		t.Errorf("Location = %q, want 'New York'", c.Location)
	}
}

func TestParseWttrJSON_EmptyConditions(t *testing.T) {
	_, err := parseWttrJSON([]byte(`{"current_condition": []}`))
	if err == nil {
		t.Error("expected error for empty conditions")
	}
}

func TestParseWttrJSON_InvalidJSON(t *testing.T) {
	_, err := parseWttrJSON([]byte("not json"))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestParseWttrJSON_ErrorResponse(t *testing.T) {
	// wttr.in text error responses won't parse as JSON
	sorry := `Sorry, we are running out of queries to the weather service at the moment.`
	_, err := parseWttrJSON([]byte(sorry))
	if err == nil {
		t.Error("expected error for text error response from wttr.in")
	}
}

func TestWttrCodeToIcon_Known(t *testing.T) {
	tests := []struct {
		code string
		want string
	}{
		{"113", "☀️"},
		{"116", "⛅"},
		{"200", "⛈️"},
		{"338", "❄️"},
	}
	for _, tt := range tests {
		got := wttrCodeToIcon(tt.code)
		if got != tt.want {
			t.Errorf("wttrCodeToIcon(%q) = %q, want %q", tt.code, got, tt.want)
		}
	}
}

func TestWttrCodeToIcon_Unknown(t *testing.T) {
	got := wttrCodeToIcon("9999")
	if got == "" {
		t.Error("expected fallback icon for unknown weather code")
	}
}

func TestLocationToPath(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"10001", "10001"},
		{"New York", "New+York"},
		{"London", "London"},
		{"40.71,-74.00", "40.71,-74.00"},
	}
	for _, tt := range tests {
		got := locationToPath(tt.input)
		if got != tt.want {
			t.Errorf("locationToPath(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
