package providers

import "testing"

const sampleOWMJSON = `{
  "weather": [{"id": 802, "description": "scattered clouds"}],
  "main": {"temp": 20.0, "feels_like": 18.0, "humidity": 60},
  "name": "Austin"
}`

func TestParseOWMJSON(t *testing.T) {
	c, err := parseOWMJSON([]byte(sampleOWMJSON))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if c.Icon != "⛅" {
		t.Errorf("Icon = %q, want ⛅", c.Icon)
	}
	if c.TempC != 20.0 {
		t.Errorf("TempC = %v, want 20.0", c.TempC)
	}
	if c.TempF != 68.0 {
		t.Errorf("TempF = %v, want 68.0", c.TempF)
	}
	if c.FeelsLikeC != 18.0 {
		t.Errorf("FeelsLikeC = %v, want 18.0", c.FeelsLikeC)
	}
	if c.Humidity != 60 {
		t.Errorf("Humidity = %v, want 60", c.Humidity)
	}
	if c.Location != "Austin" {
		t.Errorf("Location = %q, want Austin", c.Location)
	}
}

func TestParseOWMJSON_EmptyWeather(t *testing.T) {
	_, err := parseOWMJSON([]byte(`{"weather": [], "main": {}, "name": ""}`))
	if err == nil {
		t.Error("expected error for empty weather array")
	}
}

func TestParseOWMJSON_InvalidJSON(t *testing.T) {
	_, err := parseOWMJSON([]byte("not json"))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestOWMIDToIcon(t *testing.T) {
	tests := []struct {
		id   int
		want string
	}{
		{200, "⛈️"}, // thunderstorm
		{300, "🌦️"}, // drizzle
		{500, "🌦️"}, // light rain
		{502, "🌧️"}, // heavy rain
		{600, "🌨️"}, // light snow
		{602, "❄️"},  // heavy snow
		{701, "🌫️"}, // mist
		{781, "🌪️"}, // tornado
		{800, "☀️"},  // clear
		{801, "🌤️"}, // few clouds
		{802, "⛅"},  // scattered clouds
		{803, "🌥️"}, // broken clouds
		{804, "☁️"},  // overcast
	}
	for _, tt := range tests {
		got := owmIDToIcon(tt.id)
		if got != tt.want {
			t.Errorf("owmIDToIcon(%d) = %q, want %q", tt.id, got, tt.want)
		}
	}
}

func TestOWMSetLocation_CityName(t *testing.T) {
	q := make(testValues)
	owmSetLocation(q, "Austin")
	if q["q"] != "Austin" {
		t.Errorf("q = %q, want Austin", q["q"])
	}
	if q["lat"] != "" || q["lon"] != "" {
		t.Error("expected no lat/lon for city name")
	}
}

func TestOWMSetLocation_Coordinates(t *testing.T) {
	q := make(testValues)
	owmSetLocation(q, "30.26,-97.74")
	if q["lat"] != "30.26" {
		t.Errorf("lat = %q, want 30.26", q["lat"])
	}
	if q["lon"] != "-97.74" {
		t.Errorf("lon = %q, want -97.74", q["lon"])
	}
	if q["q"] != "" {
		t.Error("expected no q for coordinates")
	}
}

func TestOWMSetLocation_CityCountry(t *testing.T) {
	q := make(testValues)
	owmSetLocation(q, "London,UK")
	if q["q"] != "London,UK" {
		t.Errorf("q = %q, want London,UK", q["q"])
	}
}

// testValues is a simple map that satisfies the Set interface used by owmSetLocation.
type testValues map[string]string

func (v testValues) Set(key, val string) { v[key] = val }
