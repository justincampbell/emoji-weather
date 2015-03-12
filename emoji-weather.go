package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/justincampbell/zipcode"
	"github.com/mlbright/forecast/v2"
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
	"error":               "❗️",
}

var maxCacheAge, _ = time.ParseDuration("1h")

var coordinates string
var tempFormat string
var key string
var tmpDir string
var zipCode string

func init() {
	flag.StringVar(&coordinates, "coordinates", "", "the coordinates, expressed as latitude,longitude")
	flag.StringVar(&tempFormat, "temp", "", "display temperature in c or f")
	flag.StringVar(&key, "key", os.Getenv("FORECAST_IO_API_KEY"), "your forecast.io API key")
	flag.StringVar(&tmpDir, "tmpdir", os.TempDir(), "the directory to use to store cached responses")
	flag.StringVar(&zipCode, "zipcode", "", "a USPS ZIP Code")

	flag.Parse()
}

func main() {
	if key == "" {
		exitWith("Please provide your forecast.io API key with -key, or set FORECAST_IO_API_KEY", 1)
	}

	if tempFormat != "" && tempFormat != "c" && tempFormat != "f" {
		exitWith("The -temp argument must be 'c' or 'f'", 1)
	}

	if coordinates == "" && zipCode == "" {
		exitWith("Please provide a -zipcode or -coordinates", 1)
	}

	if zipCode != "" {
		coord, err := zipcode.Lookup(zipCode)
		check(err)
		coordinates = coord.String()
	}

	latitude, longitude := formatCoordinates(coordinates)

	cacheFilename := fmt.Sprintf("emoji-weather-%s-%s.json", latitude, longitude)
	cacheFile := path.Join(tmpDir, cacheFilename)

	var err error
	var json []byte

	if isCacheStale(cacheFile) {
		json, err = getForecast(key, latitude, longitude)
		check(err)

		err = writeCache(cacheFile, json)
		check(err)
	} else {
		json, err = ioutil.ReadFile(cacheFile)
		check(err)
	}

	fmt.Println(formatConditions(extractConditionFromJSON(json)))
}

func formatCoordinates(coordinates string) (latitude string, longitude string) {
	coordinateParts := strings.Split(coordinates, ",")

	if len(coordinateParts) != 2 {
		exitWith("You must specify latitude and longitude like so: 39.95,-75.16", 1)
	}

	latitude = roundCoordinatePart(coordinateParts[0])
	longitude = roundCoordinatePart(coordinateParts[1])

	return
}

func roundCoordinatePart(part string) string {
	float, err := strconv.ParseFloat(part, 32)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%.2f", float)
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

func formatConditions(condition string, temperature float64) (result string) {
	result, ok := conditionIcons[condition]
	if !ok {
		result = condition
	}
	if tempFormat == "c" {
		temperature = convertToCelcius(temperature)
		result = fmt.Sprintf("%s %.1f°", result, temperature)
	}
	if tempFormat == "f" {
		result = fmt.Sprintf("%s %.0f°", result, temperature)
	}
	return
}

func convertToCelcius(temperature float64) float64 {
	return (temperature - 32) * 5 / 9
}

func extractConditionFromJSON(jsonBlob []byte) (condition string, temperature float64) {
	f, err := forecast.FromJSON(jsonBlob)

	if err != nil {
		return "error", 0.0
	}

	if f.Code > 0 {
		return "error", 0.0
	}

	return f.Currently.Icon, f.Currently.Temperature
}

func exitWith(message interface{}, status int) {
	fmt.Printf("❗️\n%v\n", message)
	os.Exit(status)
}

func check(err error) {
	if err != nil {
		exitWith(err, 1)
	}
}
