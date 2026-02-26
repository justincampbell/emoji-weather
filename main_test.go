package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// mockProvider implements Provider for testing.
type mockProvider struct {
	name       string
	conditions Conditions
	err        error
	callCount  int
}

func (m *mockProvider) Name() string { return m.name }

func (m *mockProvider) Get(location string, timeout time.Duration) (Conditions, error) {
	m.callCount++
	return m.conditions, m.err
}

func TestGetCachedConditions_FetchesOnMiss(t *testing.T) {
	dir := t.TempDir()
	p := &mockProvider{
		name:       "test",
		conditions: Conditions{Icon: "☀️", TempF: 80},
	}

	got, err := getCachedConditions(p, "NYC", "u", time.Hour, 5*time.Second, dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Icon != "☀️" {
		t.Errorf("Icon = %q, want ☀️", got.Icon)
	}
	if p.callCount != 1 {
		t.Errorf("callCount = %d, want 1", p.callCount)
	}
}

func TestGetCachedConditions_UsesCache(t *testing.T) {
	dir := t.TempDir()
	p := &mockProvider{
		name:       "test",
		conditions: Conditions{Icon: "☀️", TempF: 80},
	}

	// First call populates cache.
	_, err := getCachedConditions(p, "NYC", "u", time.Hour, 5*time.Second, dir)
	if err != nil {
		t.Fatalf("first call error: %v", err)
	}

	// Second call should use cache, not call provider.
	p.conditions = Conditions{Icon: "⛈️", TempF: 55}
	got, err := getCachedConditions(p, "NYC", "u", time.Hour, 5*time.Second, dir)
	if err != nil {
		t.Fatalf("second call error: %v", err)
	}
	if got.Icon != "☀️" {
		t.Errorf("expected cached ☀️, got %q", got.Icon)
	}
	if p.callCount != 1 {
		t.Errorf("callCount = %d, want 1 (should have used cache)", p.callCount)
	}
}

func TestGetCachedConditions_RefetchesOnExpiry(t *testing.T) {
	dir := t.TempDir()
	p := &mockProvider{
		name:       "test",
		conditions: Conditions{Icon: "☀️", TempF: 80},
	}

	// Write a cache file with an old mtime.
	cachePath := filepath.Join(dir, "tmux-weather."+conditionsCacheKey("test", "NYC", "u"))
	_ = saveConditions(cachePath, p.conditions)
	oldTime := time.Now().Add(-2 * time.Hour)
	_ = os.Chtimes(cachePath, oldTime, oldTime)

	p.conditions = Conditions{Icon: "⛈️", TempF: 55}
	got, err := getCachedConditions(p, "NYC", "u", time.Hour, 5*time.Second, dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Icon != "⛈️" {
		t.Errorf("expected fresh ⛈️, got %q", got.Icon)
	}
	if p.callCount != 1 {
		t.Errorf("callCount = %d, want 1", p.callCount)
	}
}

func TestGetCachedConditions_PropagatesProviderError(t *testing.T) {
	dir := t.TempDir()
	p := &mockProvider{
		name: "test",
		err:  errors.New("network unreachable"),
	}

	_, err := getCachedConditions(p, "", "u", time.Hour, 5*time.Second, dir)
	if err == nil {
		t.Error("expected error from provider")
	}
}

func TestConditionsCacheKey_Deterministic(t *testing.T) {
	k1 := conditionsCacheKey("wttr", "NYC", "u")
	k2 := conditionsCacheKey("wttr", "NYC", "u")
	if k1 != k2 {
		t.Error("cache key not deterministic")
	}
}

func TestConditionsCacheKey_DifferentInputs(t *testing.T) {
	k1 := conditionsCacheKey("wttr", "NYC", "u")
	k2 := conditionsCacheKey("wttr", "NYC", "m")
	k3 := conditionsCacheKey("wttr", "London", "u")
	if k1 == k2 || k1 == k3 || k2 == k3 {
		t.Error("different inputs should produce different cache keys")
	}
}

func TestCacheStale_MissingFile(t *testing.T) {
	if !cacheStale("/nonexistent/path/file", time.Hour) {
		t.Error("expected stale for missing file")
	}
}

func TestCacheStale_FreshFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test")
	_ = os.WriteFile(path, []byte("data"), 0644)

	if cacheStale(path, time.Hour) {
		t.Error("expected fresh file to not be stale")
	}
}

func TestCacheStale_ExpiredFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test")
	_ = os.WriteFile(path, []byte("data"), 0644)
	oldTime := time.Now().Add(-2 * time.Hour)
	_ = os.Chtimes(path, oldTime, oldTime)

	if !cacheStale(path, time.Hour) {
		t.Error("expected old file to be stale")
	}
}

func TestNewProvider_Known(t *testing.T) {
	p, err := newProvider("wttr")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name() != "wttr" {
		t.Errorf("Name() = %q, want wttr", p.Name())
	}
}

func TestNewProvider_Unknown(t *testing.T) {
	_, err := newProvider("nonexistent")
	if err == nil {
		t.Error("expected error for unknown provider")
	}
}
