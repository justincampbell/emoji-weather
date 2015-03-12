load test_helper

@test "it's always sunny in Philadelphia with coordinates" {
  fixture emoji-weather-39.95--75.16.json

  run $ew -coordinates "39.951735,-75.158654"

  echo $output | grep "☀️"
}

@test "it's always sunny in Philadelphia with ZIP Code" {
  fixture emoji-weather-39.95--75.16.json

  run $ew -zipcode "19107"

  echo $output | grep "☀️"
}

@test "temperature display" {
  fixture emoji-weather-39.95--75.16.json

  run $ew -temp=f -zipcode "19107"

  echo $output | grep "☀️ 49°"
}

@test "show exclamation point on failure" {
  run $ew -coordinates "1000,0"

  echo $output | grep "❗️"
  [ $status -eq 0 ]
}
