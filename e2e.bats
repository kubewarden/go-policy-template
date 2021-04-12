#!/usr/bin/env bats

@test "reject because name is on deny list" {
  run policy-testdrive -p policy.wasm -r test_data/ingress.json -s '{"denied_names": ["foo", "tls-example-ingress"]}'

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # settings validation passed
  [[ "$output" == *"valid: true"* ]]

  # request rejected
  [[ "$output" == *"allowed: false"* ]]
  [[ "$output" == *"The \'tls-example-ingress\' name is on the deny list"* ]]
}

@test "accept because name is not on the deny list" {
  run policy-testdrive -p policy.wasm -r test_data/ingress.json -s '{"denied_names": ["foo"]}'
  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # settings validation passed
  [[ "$output" == *"valid: true"* ]]

  # request accepted
  [[ "$output" == *"allowed: true"* ]]
}

@test "accept because the deny list is empty" {
  run policy-testdrive -p policy.wasm -r test_data/ingress.json
  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # settings validation passed
  [[ "$output" == *"valid: true"* ]]

  # request accepted
  [[ "$output" == *"allowed: true"* ]]
}
