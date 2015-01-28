setup() {
  export FORECAST_IO_API_KEY=abc123
  ew="bin/emoji-weather -tmpdir test/tmp"
  mkdir -p test/tmp
}

teardown() {
  rm test/tmp/*.json
  echo $output
}

fixture() {
  cp test/{fixtures,tmp}/$@
  touch test/tmp/$@
}
