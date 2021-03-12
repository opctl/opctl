#!/usr/bin/env sh

echo "starting docker daemon as background process"
nohup dockerd \
  --host=unix:///var/run/docker.sock \
  --storage-driver=overlay2 &

echo "installing jq"
apk add -U jq

# dummy account for these tests so we don't hit rate limits
# it is not secret as it has no access to anything
# and is use solely for this purpose
/src/cli/opctl-linux-amd64 auth add docker.io -u 3hhyyicl1mzqsr6tggmg -p '%7Oe^4#fGGwc96rGcV&4'

test_description="opspec test-suite scenarios"

. sharness/sharness.sh

for dir in $(find /src/test-suite/**/ -type d)
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
      op: $dir
      scenario: ${name:-default}"

      case "$expect" in
        success)
          test_expect_success "$scenario_description" "
              /src/cli/opctl-linux-amd64 run --arg-file /args.yml "$dir"
          "
          ;;
        failure)
          test_expect_success "$scenario_description" "
              test_must_fail /src/cli/opctl-linux-amd64 run --arg-file /args.yml "$dir"
          "
          ;;
      esac
      done
  fi
done

test_done
