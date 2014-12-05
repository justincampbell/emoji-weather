package main

import (
	"testing"

	"github.com/bmizerany/assert"
)

func Test_formatConditions(t *testing.T) {
	assert.Equal(t, "☀️", formatConditions("clear-day"))
	assert.Equal(t, "non-existent", formatConditions("non-existent"))
}
