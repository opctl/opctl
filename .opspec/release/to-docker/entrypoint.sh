#!/bin/sh

echo 'starting docker daemon as background process'
nohup dockerd \
--host=unix:///var/run/docker.sock \
--host=tcp://0.0.0.0:2375 \
--storage-driver=overlay2 &

# poll for docker daemon up
max_retries=2
n=0
until [ $n -ge $max_retries ]
do
  docker ps > /dev/null 2>&1 && break
  n=$((n+1))
  sleep 3
done

if [ "$n" -eq "$max_retries" ]; then
  # assume failed
  cat nohup.out
  exit 1
fi

if [ $# -eq 0 ]; then
    echo 'no cmd provided, running "opctl node start"'
    opctl node start
else
  echo 'running provided cmd'
  exec "$@"
fi
