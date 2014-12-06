package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
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
	coordinates := flag.String("coordinates", "39.95,-75.1667", "the coordinates, expressed as latitude,longitude")
	flag.Parse()

	coordinateParts := strings.Split(*coordinates, ",")
	var latitude string
	var longitude string

	if len(coordinateParts) != 2 {
		exitWith("You must specify latitude and longitude like so: 39.95,-75.1667", 1)
	} else {
		latitude = coordinateParts[0]
		longitude = coordinateParts[1]
	}

	var json []byte
	var err error

	if isCacheStale(cache_file) {
		json, err = getForecast(latitude, longitude)
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

func getForecast(latitude string, longitude string) (json []byte, err error) {
	key := os.Getenv("FORECAST_IO_API_KEY")

	if key == "" {
		log.Fatal("Please set FORECAST_IO_API_KEY to your forecast.io API key")
	}

	res, err := forecast.GetResponse(key, latitude, longitude, "now", "us")
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

func exitWith(message interface{}, status int) {
	fmt.Printf("%v\n", message)
	os.Exit(status)
}

func check(err error) {
	if err != nil {
		exitWith(err, 1)
	}
}
