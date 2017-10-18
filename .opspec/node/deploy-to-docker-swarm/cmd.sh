#!/bin/sh

set -e

echo "starting docker daemon"
nohup dockerd \
  --host=unix:///var/run/docker.sock \
  --storage-driver=aufs &

# poll for docker daemon up
max_retries=6
n=0
until [ $n -ge $max_retries ]
do
  docker ps > /dev/null 2>&1 && break
  n=$((n+1))
  sleep 1
done

if [ "$n" -eq "$max_retries" ]; then
  # assume failed
  cat nohup.out
  exit 1
fi

result=$(/dockerCloudClient --user="$dockerUsername" --pass="$dockerPassword" opctl/alpha)
export "${result##* }"
sleep 1

echo creating network if not found
docker network inspect frontend || docker network create --driver overlay frontend

echo creating or updating stack
docker stack deploy -c /docker-stack.yml alpha