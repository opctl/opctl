#!/usr/bin/env sh

echo "starting docker daemon as background process"
nohup dockerd \
  --host=unix:///var/run/docker.sock \
  --storage-driver=overlay2 &

# start the node outside the $(...) capture below; a node started inside it inherits the captured stdout and the capture never ends
opctl ls "$(mktemp -d)" > /dev/null

if [[ $authAddGithub == "true" ]]; then
  opctl auth add github.com -u " " -p $githubAccessToken
fi

exec <&-

opRef="${opRef:-/test}"

echo "op: $op"
blah=$(opctl run --no-progress --arg-file /args.yml $opRef)

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
