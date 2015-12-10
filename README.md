# emoji-weather [![Build Status](https://travis-ci.org/justincampbell/emoji-weather.svg?branch=master)](https://travis-ci.org/justincampbell/emoji-weather)

The current weather expressed as Emoji.

## Installation

Download the latest binary for you platform from the [Releases page](https://github.com/justincampbell/emoji-weather/releases/latest), and place the file in your $PATH. For most systems, this will be `/usr/bin/local'.

Or, if you have a [Go development environment](https://golang.org/doc/install):

```sh
go get github.com/justincampbell/emoji-weather
```

## Usage

1. Register for a [Forecast.io Developer Account](https://developer.forecast.io/)
2. `export FORECAST_IO_API_KEY=YOUR_KEY`
3. `emoji-weather -zipcode "90210"` or `emoji-weather -coordinates "latitude,longitude"`
4. You can add `-temp=f` or `-temp=c` for temperature display

Run `emoji-weather -h` to see all available flags and their default values.
