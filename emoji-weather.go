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
	key := os.Getenv("FORECAST_IO_API_KEY")

	if key == "" {
		log.Fatal("Please set your forecast.io API key")
	}

	lat := "39.95"
	long := "-75.1667"

	var json []byte
	cache_file := path.Join(os.TempDir(), "emoji-weather.json")

	if _, err := os.Stat(cache_file); os.IsNotExist(err) {
		res, err := forecast.GetResponse(key, lat, long, "now", "us")
		if err != nil {
			panic(err)
		}

		json, err = ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(cache_file, json, 0644)
		if err != nil {
			panic(err)
		}
	} else {
		json, err = ioutil.ReadFile(cache_file)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println(formatConditions(extractConditionFromJSON(json)))
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
