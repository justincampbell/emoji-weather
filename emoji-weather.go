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
	"fog":                 "fog",
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

	var cacheFilename = fmt.Sprintf("emoji-weather-%s-%s.json", latitude, longitude)
	var cacheFile = path.Join(*tmpDir, cacheFilename)

	var json []byte
	var err error

	if isCacheStale(cacheFile) {
		json, err = getForecast(latitude, longitude)
		check(err)

		err = writeCache(cacheFile, json)
		check(err)
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

func getForecast(latitude string, longitude string) (json []byte, err error) {
	key := os.Getenv("FORECAST_IO_API_KEY")

	if key == "" {
		exitWith("Please set FORECAST_IO_API_KEY to your forecast.io API key", 1)
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
