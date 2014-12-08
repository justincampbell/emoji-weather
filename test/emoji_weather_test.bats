load test_helper

@test "it's always sunny in Philadelphia" {
  cp test/{fixtures,tmp}/emoji-weather-39.95--75.1667.json
  touch test/tmp/emoji-weather-39.95--75.1667.json

  run bin/emoji-weather -tmpdir test/tmp -coordinates "39.95,-75.1667"

  echo $output | grep "☀️"
}
