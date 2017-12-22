#!/bin/sh

echo "starting docker daemon as background process"
nohup dockerd \
--host=unix:///var/run/docker.sock \
--host=tcp://0.0.0.0:2375 \
--storage-driver=overlay2 &

echo "sleeping 2 sec to allow it to start"
sleep 2

echo "running provided cmd"
exec "$@"
