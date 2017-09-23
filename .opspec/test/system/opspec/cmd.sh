#!/usr/bin/env sh

echo "starting docker daemon as background process"
nohup dockerd > /nohup.out &

echo "installing jq"
apk add -U jq

test_description="opspec test-suite scenarios"

. sharness/sharness.sh

for dir in $(find /src/github.com/opspec-io/test-suite/scenarios/pkg/**/ -type d)
do
  if [ -f "$dir/scenarios.yml" ]
  then
    # 1) convert scenarios from yaml to json & filter to call scenarios
    # 2) run each scenario via sharness
    scenarios=$(yaml < "$dir/scenarios.yml" | jq -c '.[] | select(.call)')
    for scenario in "$scenarios"
      do
      expect=$(echo "$scenario" | jq -r '.call.expect')
      name=$(echo "$scenario" | jq -r '.name? | select(. != null)')

      # generate args.yml from scenario scope
      echo "$scenario" | jq 'select(.call.scope) | .call.scope[]' > /args.yml

      scenario_description="
      pkg: $dir
      scenario: ${name:-default}"

      case "$expect" in
        success)
          test_expect_success "$scenario_description" "
              opctl run --arg-file /args.yml "$dir"
          "
          ;;
        failure)
          test_expect_success "$scenario_description" "
              test_must_fail opctl run --arg-file /args.yml "$dir"
          "
          ;;
      esac
      done
  fi
done

test_done
