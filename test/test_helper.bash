setup() {
  mkdir -p test/tmp
}

teardown() {
  rm test/tmp/*.json
  echo $output
}
