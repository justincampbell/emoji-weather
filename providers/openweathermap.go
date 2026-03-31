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

const owmBaseURL = "https://api.openweathermap.org/data/2.5/weather"

// OpenWeatherMapProvider fetches weather from the OpenWeatherMap API.
type OpenWeatherMapProvider struct {
	apiKey string
}

func NewOpenWeatherMapProvider(apiKey string) *OpenWeatherMapProvider {
	return &OpenWeatherMapProvider{apiKey: apiKey}
}

func (p *OpenWeatherMapProvider) Name() string { return "openweathermap" }

func (p *OpenWeatherMapProvider) Get(location string, timeout time.Duration) (Conditions, error) {
	if location == "" {
		return Conditions{}, fmt.Errorf("openweathermap requires a location")
	}

	req, err := http.NewRequest("GET", owmBaseURL, nil)
	if err != nil {
		return Conditions{}, err
	}

	q := req.URL.Query()
	owmSetLocation(q, location)
	q.Set("units", "metric") // fetch in °C; convert to °F locally
	q.Set("appid", p.apiKey)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return Conditions{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return Conditions{}, fmt.Errorf("HTTP %d from openweathermap", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Conditions{}, err
	}

	return parseOWMJSON(body)
}

// owmSetLocation sets either lat/lon query params (for "lat,lon" input) or q (city name / zip).
func owmSetLocation(q interface{ Set(string, string) }, location string) {
	parts := strings.SplitN(location, ",", 2)
	if len(parts) == 2 {
		_, errLat := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		_, errLon := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if errLat == nil && errLon == nil {
			q.Set("lat", strings.TrimSpace(parts[0]))
			q.Set("lon", strings.TrimSpace(parts[1]))
			return
		}
	}
	q.Set("q", location)
}

type owmResponse struct {
	Weather []owmWeather `json:"weather"`
	Main    owmMain      `json:"main"`
	Name    string       `json:"name"`
}

type owmWeather struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type owmMain struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	Humidity  int     `json:"humidity"`
}

func parseOWMJSON(data []byte) (Conditions, error) {
	var resp owmResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return Conditions{}, fmt.Errorf("openweathermap parse error: %w", err)
	}

	if len(resp.Weather) == 0 {
		return Conditions{}, fmt.Errorf("openweathermap returned no weather data")
	}

	tempC := resp.Main.Temp
	feelsC := resp.Main.FeelsLike
	tempF := tempC*9/5 + 32
	feelsF := feelsC*9/5 + 32

	w := resp.Weather[0]
	desc := strings.Title(w.Description)

	return Conditions{
		Icon:        owmIDToIcon(w.ID),
		Description: desc,
		TempF:       tempF,
		TempC:       tempC,
		FeelsLikeF:  feelsF,
		FeelsLikeC:  feelsC,
		Humidity:    resp.Main.Humidity,
		Location:    resp.Name,
	}, nil
}

// owmIDToIcon maps OpenWeatherMap condition IDs to emoji.
// See https://openweathermap.org/weather-conditions
func owmIDToIcon(id int) string {
	switch {
	case id >= 200 && id < 300:
		return "⛈️" // Thunderstorm
	case id >= 300 && id < 400:
		return "🌦️" // Drizzle
	case id >= 500 && id < 600:
		switch id {
		case 500, 520:
			return "🌦️" // Light rain / light shower
		default:
			return "🌧️"
		}
	case id >= 600 && id < 700:
		switch id {
		case 602, 621, 622:
			return "❄️" // Heavy snow
		default:
			return "🌨️"
		}
	case id >= 700 && id < 800:
		switch id {
		case 771:
			return "💨" // Squall
		case 781:
			return "🌪️" // Tornado
		default:
			return "🌫️"
		}
	case id == 800:
		return "☀️" // Clear sky
	case id == 801:
		return "🌤️" // Few clouds
	case id == 802:
		return "⛅" // Scattered clouds
	case id == 803:
		return "🌥️" // Broken clouds
	case id == 804:
		return "☁️" // Overcast
	}
	return "🌡️"
}
