package main

import (
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/justincampbell/emoji-weather/providers"
)

var version = "dev"

var (
	locationFlag  = flag.String("location", "", "location (zip, city, lat,lon); defaults to auto-detect by IP")
	formatFlag    = flag.String("format", "%c %t", "format string: %c=icon %C=description %t=temp %f=feels-like %h=humidity %l=location")
	ttlFlag       = flag.Duration("ttl", 30*time.Minute, "cache TTL (e.g. 5m, 1h)")
	unitsFlag     = flag.String("units", "f", "units: f (°F) or c (°C)")
	timeoutFlag   = flag.Duration("timeout", 5*time.Second, "HTTP request timeout")
	errorIconFlag = flag.String("error-icon", "❗", "output on error (should not produce noise in tmux)")
	cacheDirFlag  = flag.String("cache-dir", "", "cache directory (defaults to system temp dir)")
	providerFlag  = flag.String("provider", "wttr", "weather provider: wttr, openweathermap")
	apiKeyFlag    = flag.String("api-key", "", "API key for weather provider (openweathermap reads ~/.openweather by default)")
	verboseFlag   = flag.Bool("verbose", false, "print errors to stderr instead of suppressing them")
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

	apiKey := *apiKeyFlag
	if apiKey == "" && *providerFlag == "openweathermap" {
		if home, err := os.UserHomeDir(); err == nil {
			if data, err := os.ReadFile(filepath.Join(home, ".openweather")); err == nil {
				apiKey = strings.TrimSpace(string(data))
			}
		}
	}

	location := *locationFlag
	if location == "" {
		loc, err := detectLocation(*timeoutFlag)
		if err != nil {
			if *verboseFlag {
				fmt.Fprintln(os.Stderr, "error detecting location:", err)
			}
			fmt.Print(*errorIconFlag)
			os.Exit(0)
		}
		location = loc
	}

	provider, err := newProvider(*providerFlag, apiKey)
	if err != nil {
		if *verboseFlag {
			fmt.Fprintln(os.Stderr, "error:", err)
		}
		fmt.Print(*errorIconFlag)
		os.Exit(0)
	}

	conditions, err := getCachedConditions(provider, location, *unitsFlag, *ttlFlag, *timeoutFlag, cacheDir)
	if err != nil {
		if *verboseFlag {
			fmt.Fprintln(os.Stderr, "error:", err)
		}
		fmt.Print(*errorIconFlag)
		os.Exit(0)
	}

	fmt.Print(Format(*formatFlag, conditions, *unitsFlag))
}

// newProvider returns the named Provider, or an error if unknown.
func newProvider(name, apiKey string) (providers.Provider, error) {
	switch name {
	case "wttr":
		return providers.NewWttrProvider(version), nil
	case "openweathermap":
		if apiKey == "" {
			return nil, fmt.Errorf("openweathermap requires an API key (set --api-key or place key in ~/.openweather)")
		}
		return providers.NewOpenWeatherMapProvider(apiKey), nil
	default:
		return nil, fmt.Errorf("unknown provider %q", name)
	}
}

func getCachedConditions(p providers.Provider, location, units string, ttl, timeout time.Duration, cacheDir string) (providers.Conditions, error) {
	cachePath := filepath.Join(cacheDir, "emoji-weather."+conditionsCacheKey(p.Name(), location, units))

	if !cacheStale(cachePath, ttl) {
		if c, err := loadConditions(cachePath); err == nil {
			return c, nil
		}
	}

	conditions, err := p.Get(location, timeout)
	if err != nil {
		return providers.Conditions{}, err
	}

	_ = saveConditions(cachePath, conditions)
	return conditions, nil
}

func loadConditions(path string) (providers.Conditions, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return providers.Conditions{}, err
	}
	var c providers.Conditions
	return c, json.Unmarshal(data, &c)
}

func saveConditions(path string, c providers.Conditions) error {
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

// detectLocation returns the user's approximate coordinates via ipinfo.io.
func detectLocation(timeout time.Duration) (string, error) {
	client := &http.Client{Timeout: timeout}
	resp, err := client.Get("https://ipinfo.io/json")
	if err != nil {
		return "", fmt.Errorf("IP location detection failed: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	var info struct {
		Loc string `json:"loc"` // "lat,lon"
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return "", fmt.Errorf("IP location detection: bad response: %w", err)
	}
	if info.Loc == "" {
		return "", fmt.Errorf("IP location detection returned empty location")
	}
	return info.Loc, nil
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
