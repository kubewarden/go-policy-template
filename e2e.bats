#!/usr/bin/env bats

@test "reject because name is on deny list" {
  run kwctl run policy.wasm -r test_data/ingress.json --settings-json '{"denied_names": ["foo", "tls-example-ingress"]}'

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request rejected
  [ "$status" -eq 0 ]
  [ $(expr "$output" : '.*allowed.*false') -ne 0 ]
  [ $(expr "$output" : ".*The 'tls-example-ingress' name is on the deny list.*") -ne 0 ]
}

@test "accept because name is not on the deny list" {
  run kwctl run policy.wasm -r test_data/ingress.json --settings-json '{"denied_names": ["foo"]}'
  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request accepted
  [ "$status" -eq 0 ]
  [ $(expr "$output" : '.*allowed.*true') -ne 0 ]
}

@test "accept because the deny list is empty" {
  run kwctl run policy.wasm -r test_data/ingress.json
  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request accepted
  [ "$status" -eq 0 ]
  [ $(expr "$output" : '.*allowed.*true') -ne 0 ]
}
