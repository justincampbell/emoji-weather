load test_helper

@test "it's always sunny in Philadelphia" {
  run bin/emoji-weather

  echo $output | grep "☀️"
}
