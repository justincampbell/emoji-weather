package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_formatConditions(t *testing.T) {
	assert.Equal(t, "☀️", formatConditions("clear-day"))
	assert.Equal(t, "non-existent", formatConditions("non-existent"))
}

func Test_extractCloudyConditionFromJSON(t *testing.T) {
	json := `{ "currently": { "icon": "cloudy" } }`
	jsonBlob := []byte(json)

	assert.Equal(t, "cloudy", extractConditionFromJSON(jsonBlob))
}

func Test_extractClearConditionFromJSON(t *testing.T) {
	json := `{ "currently": { "icon": "clear-day" } }`
	jsonBlob := []byte(json)

	assert.Equal(t, "clear-day", extractConditionFromJSON(jsonBlob))
}
