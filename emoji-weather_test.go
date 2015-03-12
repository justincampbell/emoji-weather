package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_formatConditions(t *testing.T) {
	assert.Equal(t, "☀️", formatConditions("clear-day", 0.0))
	assert.Equal(t, "non-existent", formatConditions("non-existent", 0.0))
}

func Test_extractCloudyConditionFromJSON(t *testing.T) {
	json := `{ "currently": { "icon": "cloudy" } }`
	jsonBlob := []byte(json)

	icon, temperature := extractConditionFromJSON(jsonBlob)
	assert.Equal(t, "cloudy", icon)
	assert.Equal(t, 0.0, temperature)
}

func Test_convertToCelcius(t *testing.T) {
	assert.InDelta(t, 100.0, convertToCelcius(212.0), 0.1)
	assert.InDelta(t, 21.1, convertToCelcius(70.0), 0.1)
	assert.InDelta(t, -5.5, convertToCelcius(22.1), 0.1)
}

func Test_formatCoordinates(t *testing.T) {
	coordinates := "12.34567,-65.43210"
	lat, long := formatCoordinates(coordinates)
	assert.Equal(t, "12.35", lat)
	assert.Equal(t, "-65.43", long)
}
