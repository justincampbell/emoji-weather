package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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
	"fog":                 "🌁",
	"partly-cloudy-day":   "⛅️",
	"partly-cloudy-night": "🌙",
	"rain":                "☔️",
	"sleet":               "❄️ ☔️",
	"snow":                "❄️",
	"wind":                "🍃",
}

var maxCacheAge, _ = time.ParseDuration("1h")

func main() {
	coordinates := flag.String("coordinates", "39.95,-75.1667", "the coordinates, expressed as latitude,longitude")
	tmpDir := flag.String("tmpdir", os.TempDir(), "the directory to use to store cached responses")
	key := flag.String("key", os.Getenv("FORECAST_IO_API_KEY"), "your forecast.io API key")

	if *key == "" {
		exitWith("Please provide your forecast.io API key with -key, or set FORECAST_IO_API_KEY", 1)
	}

	flag.Parse()

	coordinateParts := strings.Split(*coordinates, ",")

	if len(coordinateParts) != 2 {
		exitWith("You must specify latitude and longitude like so: 39.95,-75.1667", 1)
	}

	latitude := coordinateParts[0]
	longitude := coordinateParts[1]

	cacheFilename := fmt.Sprintf("emoji-weather-%s-%s.json", latitude, longitude)
	cacheFile := path.Join(*tmpDir, cacheFilename)

	var json []byte
	var err error

	if isCacheStale(cacheFile) {
		json, err = getForecast(*key, latitude, longitude)
		if err == nil {
			err = writeCache(cacheFile, json)
			check(err)
		}
	} else {
		json, err = ioutil.ReadFile(cacheFile)
		check(err)
	}

	fmt.Println(formatConditions(extractConditionFromJSON(json)))
}

func isCacheStale(cacheFile string) bool {
	stat, err := os.Stat(cacheFile)

	return os.IsNotExist(err) || time.Since(stat.ModTime()) > maxCacheAge
}

func getForecast(key string, latitude string, longitude string) (json []byte, err error) {
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

func writeCache(cacheFile string, json []byte) (err error) {
	return ioutil.WriteFile(cacheFile, json, 0644)
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
