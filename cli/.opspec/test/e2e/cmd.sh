#!/usr/bin/env sh

echo "starting docker daemon as background process"
nohup dockerd \
  --host=unix:///var/run/docker.sock \
  --storage-driver=overlay2 &

# dummy account for these tests so we don't hit rate limits
# it is not secret as it has no access to anything
# and is use solely for this purpose
opctl auth add docker.io -u 3hhyyicl1mzqsr6tggmg -p '%7Oe^4#fGGwc96rGcV&4'

exec <&-

echo "op: $op"
blah=$(opctl run --no-progress --arg-file /args.yml /test)

case "$expect" in
  success)
    if [ $? -eq 0 ]; then
      echo "expected $expect and got success"
      exit 0
    else
      echo "expected $expect but got failure"
      exit 1
    fi
    ;;
  failure)
    if [ $? -eq 0 ]; then
      echo "expected $expect but got success"
      exit 1
    else
      echo "expected $expect and got failure"
      exit 0
    fi
    ;;
esac
