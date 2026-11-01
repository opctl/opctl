#!/usr/bin/env sh

if [ "${cacheProof:-false}" = true ]; then
  export OPCTL_DATA_DIR=/tmp/opctl-cache-proof
fi

echo "starting docker daemon as background process"
nohup dockerd \
  --host=unix:///var/run/docker.sock \
  --storage-driver=overlay2 &

# dummy account for these tests so we don't hit rate limits
# it is not secret as it has no access to anything
# and is use solely for this purpose
opctl auth add docker.io -u 3hhyyicl1mzqsr6tggmg -p '%7Oe^4#fGGwc96rGcV&4'

if [ "$authAddGithub" = "true" ]; then
  opctl auth add github.com -u " " -p $githubAccessToken
fi

exec <&-

if [ -n "$createGo" ]; then
  createParent="${createParent:-/tmp}"
  opctl op create --go -p "$createParent" -d 'Greets the caller.' "$createGo" || exit 1
  # The e2e CLI is a development build, so the scaffold needs the in-repo SDK.
  printf '\nreplace github.com/opctl/opctl/sdks/go/opspec/sdk => %s\n' "${sdkPath:-/repo/sdks/go/opspec/sdk}" >> "$createParent/$createGo/go.mod"
  opRef="$createParent/$createGo"
fi

opRef="${opRef:-/test}"

echo "op: $op"
if [ "${cacheProof:-false}" = true ]; then
  exec "${assertCachePath:-/assert-cache}" "$opRef" /args.yml
fi
opctl run --no-progress --arg-file /args.yml "$opRef"
runStatus=$?

case "$expect" in
  success)
    if [ "$runStatus" -eq 0 ]; then
      echo "expected $expect and got success"
      expectedOutputsPath="${expectedOutputsPath:-/expected-outputs.json}"
      if [ -s "$expectedOutputsPath" ]; then
        "${assertOutputsPath:-/assert-outputs}" "$expectedOutputsPath" "$opRef" || exit 1
      fi
      exit 0
    else
      echo "expected $expect but got failure"
      exit 1
    fi
    ;;
  failure)
    if [ "$runStatus" -eq 0 ]; then
      echo "expected $expect but got success"
      exit 1
    else
      echo "expected $expect and got failure"
      exit 0
    fi
    ;;
esac

echo "unknown expected outcome: $expect"
exit 1
