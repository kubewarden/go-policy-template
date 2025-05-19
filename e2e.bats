#!/usr/bin/env bats

@test "accept when no settings are provided" {
  run kwctl run -r test_data/pod.json policy.wasm

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request is accepted
  [ $(expr "$output" : '.*"allowed":true.*') -ne 0 ]
}

@test "accept because label is satisfying a constraint" {
  run kwctl run annotated-policy.wasm \
    -r test_data/pod.json \
    --settings-json '{"constrained_labels": {"cc-center": "\\d+"}}'

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  [ "$status" -eq 0 ]
  [ $(expr "$output" : '.*allowed.*true') -ne 0 ]
}

@test "accept labels are not on deny list" {
  run kwctl run \
    -r test_data/pod.json \
    --settings-json '{"denied_labels": ["foo", "bar"]}' \
    policy.wasm

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  [ $(expr "$output" : '.*"allowed":true.*') -ne 0 ]
}

@test "reject because label is on deny list" {
  run kwctl run annotated-policy.wasm \
    -r test_data/pod.json \
    --settings-json '{"denied_labels": ["foo", "owner"]}'

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  [ "$status" -eq 0 ]
  [ $(expr "$output" : '.*allowed.*false') -ne 0 ]
  [ $(expr "$output" : ".*label owner is on the deny list.*") -ne 0 ]
}

@test "reject because label is not satisfying a constraint" {
  run kwctl run annotated-policy.wasm \
    -r test_data/pod.json \
    --settings-json '{"constrained_labels": {"cc-center": "team-\\d+"}}'

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  [ "$status" -eq 0 ]
  [ $(expr "$output" : '.*allowed.*false') -ne 0 ]
  [ $(expr "$output" : ".*the value of cc-center doesn't pass user-defined constraint.*") -ne 0 ]
}

@test "reject because constrained label is missing" {
  run kwctl run annotated-policy.wasm \
    -r test_data/pod.json \
    --settings-json '{"constrained_labels": {"organization": "\\d+"}}'

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  [ "$status" -eq 0 ]
  [ $(expr "$output" : '.*allowed.*false') -ne 0 ]
  [ $(expr "$output" : ".*Constrained label organization not found inside of Pod.*") -ne 0 ]
}

@test "fail settings validation because of conflicting labels" {
  run kwctl run \
    -r test_data/pod.json \
    --settings-json '{"denied_labels": ["foo", "cc-center"], "constrained_labels": {"cc-center": "^cc-\\d+$"}}' \
    policy.wasm

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # settings validation failed
  [ $(expr "$output" : ".*Provided settings are not valid: these labels cannot be constrained and denied at the same time: Set{cc-center}.*") -ne 0 ]
}

@test "fail settings validation because of invalid constraint" {
  run kwctl run \
    -r test_data/pod.json \
    --settings-json '{"constrained_labels": {"cc-center": "^cc-[12$"}}' \
    policy.wasm

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # settings validation failed
  [ $(expr "$output" : ".*Provided settings are not valid: error parsing regexp.*") -ne 0 ]
}