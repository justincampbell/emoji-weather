package providers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const wttrBaseURL = "https://wttr.in"

// WttrProvider fetches weather from wttr.in using its JSON API.
type WttrProvider struct{ version string }

func NewWttrProvider(version string) *WttrProvider {
	return &WttrProvider{version: version}
}

func (p *WttrProvider) Name() string { return "wttr" }

func (p *WttrProvider) Get(location string, timeout time.Duration) (Conditions, error) {
	rawURL := wttrBaseURL + "/"
	if location != "" {
		rawURL += locationToPath(location)
	}

	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return Conditions{}, err
	}

	q := req.URL.Query()
	q.Set("format", "j1")
	req.URL.RawQuery = q.Encode()
	req.Header.Set("User-Agent", "emoji-weather/"+p.version)

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return Conditions{}, err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != 200 {
		return Conditions{}, fmt.Errorf("HTTP %d from wttr.in", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Conditions{}, err
	}

	return parseWttrJSON(body)
}

// locationToPath converts a human-readable location to a wttr.in URL path segment.
// Spaces are replaced with + per wttr.in conventions.
func locationToPath(location string) string {
	return strings.ReplaceAll(location, " ", "+")
}

type wttrResponse struct {
	CurrentCondition []wttrCurrentCondition `json:"current_condition"`
	NearestArea      []wttrNearestArea      `json:"nearest_area"`
}

type wttrCurrentCondition struct {
	TempF       string         `json:"temp_F"`
	TempC       string         `json:"temp_C"`
	FeelsLikeF  string         `json:"FeelsLikeF"`
	FeelsLikeC  string         `json:"FeelsLikeC"`
	Humidity    string         `json:"humidity"`
	WeatherCode string         `json:"weatherCode"`
	WeatherDesc []wttrDescItem `json:"weatherDesc"`
}

type wttrDescItem struct {
	Value string `json:"value"`
}

type wttrNearestArea struct {
	AreaName []wttrDescItem `json:"areaName"`
}

func parseWttrJSON(data []byte) (Conditions, error) {
	var resp wttrResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return Conditions{}, fmt.Errorf("wttr.in parse error: %w", err)
	}

	if len(resp.CurrentCondition) == 0 {
		return Conditions{}, fmt.Errorf("wttr.in returned no current conditions")
	}

	cur := resp.CurrentCondition[0]

	tempF, _ := strconv.ParseFloat(cur.TempF, 64)
	tempC, _ := strconv.ParseFloat(cur.TempC, 64)
	feelsF, _ := strconv.ParseFloat(cur.FeelsLikeF, 64)
	feelsC, _ := strconv.ParseFloat(cur.FeelsLikeC, 64)
	humidity, _ := strconv.Atoi(cur.Humidity)

	var desc string
	if len(cur.WeatherDesc) > 0 {
		desc = cur.WeatherDesc[0].Value
	}

	var loc string
	if len(resp.NearestArea) > 0 && len(resp.NearestArea[0].AreaName) > 0 {
		loc = resp.NearestArea[0].AreaName[0].Value
	}

	return Conditions{
		Icon:        wttrCodeToIcon(cur.WeatherCode),
		Description: desc,
		TempF:       tempF,
		TempC:       tempC,
		FeelsLikeF:  feelsF,
		FeelsLikeC:  feelsC,
		Humidity:    humidity,
		Location:    loc,
	}, nil
}

// wttrCodeToIcon maps wttr.in weather codes to emoji.
// See https://www.worldweatheronline.com/feed/wwo-weather-code.txt
var wttrIcons = map[string]string{
	"113": "☀️",  // Clear/Sunny
	"116": "⛅",  // Partly cloudy
	"119": "☁️",  // Cloudy
	"122": "☁️",  // Overcast
	"143": "🌫️", // Mist
	"176": "🌦️", // Patchy rain
	"179": "🌨️", // Patchy snow
	"182": "🌧️", // Patchy sleet
	"185": "🌧️", // Patchy freezing drizzle
	"200": "⛈️", // Thundery outbreaks
	"227": "🌨️", // Blowing snow
	"230": "❄️",  // Blizzard
	"248": "🌫️", // Fog
	"260": "🌫️", // Freezing fog
	"263": "🌦️", // Patchy light drizzle
	"266": "🌧️", // Light drizzle
	"281": "🌧️", // Freezing drizzle
	"284": "🌧️", // Heavy freezing drizzle
	"293": "🌦️", // Patchy light rain
	"296": "🌧️", // Light rain
	"299": "🌧️", // Moderate rain at times
	"302": "🌧️", // Moderate rain
	"305": "🌧️", // Heavy rain at times
	"308": "🌧️", // Heavy rain
	"311": "🌧️", // Light freezing rain
	"314": "🌧️", // Moderate or heavy freezing rain
	"317": "🌧️", // Light sleet
	"320": "🌧️", // Moderate or heavy sleet
	"323": "🌨️", // Patchy light snow
	"326": "🌨️", // Light snow
	"329": "🌨️", // Patchy moderate snow
	"332": "🌨️", // Moderate snow
	"335": "🌨️", // Patchy heavy snow
	"338": "❄️",  // Heavy snow
	"350": "🌧️", // Ice pellets
	"353": "🌦️", // Light rain shower
	"356": "🌧️", // Moderate or heavy rain shower
	"359": "🌧️", // Torrential rain shower
	"362": "🌧️", // Light sleet showers
	"365": "🌧️", // Moderate or heavy sleet showers
	"368": "🌨️", // Light snow showers
	"371": "🌨️", // Moderate or heavy snow showers
	"374": "🌧️", // Light showers of ice pellets
	"377": "🌧️", // Moderate or heavy showers of ice pellets
	"386": "⛈️", // Patchy light rain with thunder
	"389": "⛈️", // Moderate or heavy rain with thunder
	"392": "⛈️", // Patchy light snow with thunder
	"395": "⛈️", // Moderate or heavy snow with thunder
}

func wttrCodeToIcon(code string) string {
	if icon, ok := wttrIcons[code]; ok {
		return icon
	}
	return "🌡️"
}
