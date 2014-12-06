package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

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

var maxCacheAge, _ = time.ParseDuration("1h")
var cache_file = path.Join(os.TempDir(), "emoji-weather.json")

func main() {
	var json []byte
	var err error

	if isCacheStale(cache_file) {
		json, err = getForecast()
		check(err)

		err = writeCache(cache_file, json)
		check(err)
	} else {
		json, err = ioutil.ReadFile(cache_file)
		check(err)
	}

	fmt.Println(formatConditions(extractConditionFromJSON(json)))
}

func isCacheStale(cache_file string) bool {
	stat, err := os.Stat(cache_file)

	return os.IsNotExist(err) || time.Since(stat.ModTime()) > maxCacheAge
}

func getForecast() (json []byte, err error) {
	key := os.Getenv("FORECAST_IO_API_KEY")

	if key == "" {
		log.Fatal("Please set FORECAST_IO_API_KEY to your forecast.io API key")
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

	return json, nil
}

func writeCache(cache_file string, json []byte) (err error) {
	return ioutil.WriteFile(cache_file, json, 0644)
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

func check(err error) {
	if err != nil {
		panic(err)
	}
}
