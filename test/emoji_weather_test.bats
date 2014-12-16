load test_helper

@test "it's always sunny in Philadelphia" {
  fixture emoji-weather-39.95--75.1667.json

  run $ew -coordinates "39.95,-75.1667"

  echo $output | grep "☀️"
}
