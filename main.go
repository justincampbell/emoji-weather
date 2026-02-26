package main

import (
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

var version = "dev"

var (
	locationFlag  = flag.String("location", "", "location (zip, city, lat,lon); defaults to auto-detect by IP")
	formatFlag    = flag.String("format", "%c %t", "format string: %c=icon %C=description %t=temp %f=feels-like %h=humidity %l=location")
	ttlFlag       = flag.Duration("ttl", 30*time.Minute, "cache TTL (e.g. 5m, 1h)")
	unitsFlag     = flag.String("units", "u", "units: u (US/imperial °F), m (metric °C)")
	timeoutFlag   = flag.Duration("timeout", 10*time.Second, "HTTP request timeout")
	errorIconFlag = flag.String("error-icon", "❗", "output on error (should not produce noise in tmux)")
	cacheDirFlag  = flag.String("cache-dir", "", "cache directory (defaults to system temp dir)")
	providerFlag  = flag.String("provider", "wttr", "weather provider (currently: wttr)")
	versionFlag   = flag.Bool("version", false, "print version and exit")
)

func main() {
	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		return
	}

	cacheDir := *cacheDirFlag
	if cacheDir == "" {
		cacheDir = os.TempDir()
	}

	provider, err := newProvider(*providerFlag)
	if err != nil {
		fmt.Print(*errorIconFlag)
		os.Exit(0)
	}

	conditions, err := getCachedConditions(provider, *locationFlag, *unitsFlag, *ttlFlag, *timeoutFlag, cacheDir)
	if err != nil {
		fmt.Print(*errorIconFlag)
		os.Exit(0)
	}

	fmt.Print(Format(*formatFlag, conditions, *unitsFlag))
}

// newProvider returns the named Provider, or an error if unknown.
func newProvider(name string) (Provider, error) {
	switch name {
	case "wttr":
		return NewWttrProvider(), nil
	default:
		return nil, fmt.Errorf("unknown provider %q", name)
	}
}

func getCachedConditions(p Provider, location, units string, ttl, timeout time.Duration, cacheDir string) (Conditions, error) {
	cachePath := filepath.Join(cacheDir, "tmux-weather."+conditionsCacheKey(p.Name(), location, units))

	if !cacheStale(cachePath, ttl) {
		if c, err := loadConditions(cachePath); err == nil {
			return c, nil
		}
	}

	conditions, err := p.Get(location, timeout)
	if err != nil {
		return Conditions{}, err
	}

	_ = saveConditions(cachePath, conditions)
	return conditions, nil
}

func loadConditions(path string) (Conditions, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Conditions{}, err
	}
	var c Conditions
	return c, json.Unmarshal(data, &c)
}

func saveConditions(path string, c Conditions) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func cacheStale(path string, ttl time.Duration) bool {
	stat, err := os.Stat(path)
	return os.IsNotExist(err) || time.Since(stat.ModTime()) > ttl
}

// conditionsCacheKey hashes the provider name, location, and units.
// Format is intentionally excluded so changing --format never requires a new fetch.
func conditionsCacheKey(parts ...string) string {
	h := md5.New()
	for _, p := range parts {
		_, _ = io.WriteString(h, p)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
