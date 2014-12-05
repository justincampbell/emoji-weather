package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/justincampbell/forecast/v2"
)

var conditionIcons = map[string]string{
	"clear-day":           "☀️",
	"clear-night":         "🌙",
	"cloudy":              "☁️",
	"fog":                 "fog",
	"partly-cloudy-day":   "⛅️",
	"partly-cloudy-night": "🌙",
	"rain":                "☔️",
	"sleet":               "❄️ ☔️",
	"snow":                "❄️",
	"wind":                "🍃",
}

func main() {
	var json []byte
	var err error
	cache_file := path.Join(os.TempDir(), "emoji-weather.json")

	if isCacheStale(cache_file) {
		json, err = writeCache(cache_file)
	} else {
		json, err = ioutil.ReadFile(cache_file)
	}

	if err != nil {
		panic(err)
	}

	fmt.Println(formatConditions(extractConditionFromJSON(json)))
}

func isCacheStale(cache_file string) bool {
	_, err := os.Stat(cache_file)

	return os.IsNotExist(err)
}

func writeCache(cache_file string) (json []byte, err error) {
	key := os.Getenv("FORECAST_IO_API_KEY")

	if key == "" {
		log.Fatal("Please set your forecast.io API key")
	}

	lat := "39.95"
	long := "-75.1667"

	res, err := forecast.GetResponse(key, lat, long, "now", "us")
	if err != nil {
		return nil, err
	}

	json, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(cache_file, json, 0644)
	if err != nil {
		return nil, err
	}

	return json, nil
}

func formatConditions(condition string) (icon string) {
	icon, ok := conditionIcons[condition]
	if !ok {
		icon = condition
	}
	return
}

func extractConditionFromJSON(jsonBlob []byte) (condition string) {
	f, err := forecast.FromJSON(jsonBlob)
	if err != nil {
		return "❗️"
	}

	return f.Currently.Icon
}
