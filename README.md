# emoji-weather

Emoji weather for your terminal and tmux status bar.

## Install

```
brew install justincampbell/tap/emoji-weather
```

## Usage

```
emoji-weather
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--location` | auto-detect | Location: zip, city, or lat,lon |
| `--format` | `%c %t` | Format string (see below) |
| `--units` | `f` | `f` for °F, `c` for °C |
| `--provider` | `wttr` | `wttr` or `openweathermap` |
| `--api-key` | | API key (openweathermap reads `~/.openweather` by default) |
| `--ttl` | `30m` | Cache TTL |
| `--timeout` | `5s` | HTTP timeout |
| `--verbose` | | Print errors to stderr |
| `--error-icon` | `❗` | Output on error |

### Format codes

| Code | Description |
|------|-------------|
| `%c` | Weather icon (emoji) |
| `%C` | Description (e.g. "Partly cloudy") |
| `%t` | Temperature |
| `%f` | Feels-like temperature |
| `%h` | Humidity |
| `%l` | Location name |
| `%%` | Literal `%` |

## tmux

Add to your `tmux.conf`:

```
set -g status-right "#{weather_status} %H:%M"
```

### tmux options

```
set -g @weather_location  "10001"
set -g @weather_format    "%c %t"
set -g @weather_ttl       "30m"
set -g @weather_units     "f"
set -g @weather_provider  "wttr"
```
