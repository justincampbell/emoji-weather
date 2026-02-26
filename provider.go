package main

import "time"

// Conditions holds weather data returned by any provider.
type Conditions struct {
	Icon        string  `json:"icon"`         // weather emoji
	Description string  `json:"description"`  // e.g. "Partly cloudy"
	TempF       float64 `json:"temp_f"`
	TempC       float64 `json:"temp_c"`
	FeelsLikeF  float64 `json:"feels_like_f"`
	FeelsLikeC  float64 `json:"feels_like_c"`
	Humidity    int     `json:"humidity"`
	Location    string  `json:"location"`
}

// Provider fetches weather conditions for a given location.
type Provider interface {
	Name() string
	Get(location string, timeout time.Duration) (Conditions, error)
}
