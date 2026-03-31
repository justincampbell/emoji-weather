#!/usr/bin/env bash
# emoji-weather.tmux — tmux plugin entry point
#
# Replaces #{weather_status} in status-right, status-left, and status-format
# with a #(emoji-weather ...) shell command invocation.
#
# Configuration (set in tmux.conf):
#   set -g @weather_location  "10001"       # zip, city, lat,lon; default: auto
#   set -g @weather_format    "%c %t"       # format string; default: "%c %t"
#   set -g @weather_ttl       "30m"         # cache TTL; default: 30m
#   set -g @weather_units     "f"           # f=Fahrenheit, c=Celsius; default: f
#   set -g @weather_error_icon "❗"         # output on error; default: ❗
#   set -g @weather_provider  "wttr"        # provider; default: wttr
#
# Usage:
#   set -g status-right "#{weather_status} %H:%M"

CURRENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

find_binary() {
  # Prefer a binary built alongside the plugin (development / bundled installs).
  if [ -x "$CURRENT_DIR/bin/emoji-weather" ]; then
    echo "$CURRENT_DIR/bin/emoji-weather"
    return
  fi
  # Fall back to whatever is on PATH (e.g. installed via Homebrew).
  if command -v emoji-weather >/dev/null 2>&1; then
    command -v emoji-weather
    return
  fi
  echo ""
}

get_option() {
  local option="$1"
  local default="$2"
  local value
  value=$(tmux show-option -gqv "$option")
  echo "${value:-$default}"
}

build_command() {
  local binary="$1"
  local location format ttl units error_icon provider

  location=$(get_option "@weather_location" "")
  format=$(get_option "@weather_format" "")
  ttl=$(get_option "@weather_ttl" "")
  units=$(get_option "@weather_units" "")
  error_icon=$(get_option "@weather_error_icon" "")
  provider=$(get_option "@weather_provider" "")

  local cmd="$binary"
  [ -n "$location" ]   && cmd="$cmd --location $(printf '%q' "$location")"
  [ -n "$format" ]     && cmd="$cmd --format $(printf '%q' "$format")"
  [ -n "$ttl" ]        && cmd="$cmd --ttl $ttl"
  [ -n "$units" ]      && cmd="$cmd --units $units"
  [ -n "$error_icon" ] && cmd="$cmd --error-icon $(printf '%q' "$error_icon")"
  [ -n "$provider" ]   && cmd="$cmd --provider $provider"

  echo "$cmd"
}

update_option() {
  local option="$1"
  local interpolation="$2"
  local placeholder="#{weather_status}"
  local value
  value=$(tmux show-option -gqv "$option")
  if [[ "$value" == *"$placeholder"* ]]; then
    tmux set-option -gq "$option" "${value//$placeholder/$interpolation}"
  fi
}

main() {
  local binary
  binary=$(find_binary)

  if [ -z "$binary" ]; then
    tmux display-message "emoji-weather: binary not found — install via: brew install justincampbell/tap/emoji-weather"
    return 1
  fi

  local cmd
  cmd=$(build_command "$binary")
  local interpolation="#($cmd)"

  update_option "status-right" "$interpolation"
  update_option "status-left" "$interpolation"
  update_option "status-format[0]" "$interpolation"
  update_option "status-format[1]" "$interpolation"
}

main
