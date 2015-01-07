load test_helper

@test "it's always sunny in Philadelphia" {
  fixture emoji-weather-39.95--75.1667.json

  run $ew -coordinates "39.95,-75.1667"

  echo $output | grep "☀️"
}

@test "show exclamation point on failure" {
  run $ew -coordinates "0,0"

  echo $output | grep "❗️"
  [ $status -eq 0 ]
}
