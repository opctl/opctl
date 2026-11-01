#!/bin/sh
set -eu
umask 077
owner=/opctl/mod
[ ! -d "$owner" ] && owner=/opctl/program
ids=$(stat -c '%u:%g' "$owner") || { echo 'Go bootstrap cache must be shared at the same host path with the container daemon' >&2; exit 1; }
uid=${ids%:*}
gid=${ids#*:}
command -v setpriv >/dev/null || { echo 'Go bootstrap requires setpriv to drop privileges' >&2; exit 1; }
scratch=$(mktemp -d /tmp/opctl-bootstrap.XXXXXX)
trap 'rm -rf "$scratch"' EXIT HUP INT TERM
mkdir "$scratch/home" "$scratch/tmp" "$scratch/gopath"
chown -R "$uid:$gid" "$scratch"
if [ -d /opctl/build ]; then chown "$uid:$gid" /opctl/build; fi
export HOME="$scratch/home" TMPDIR="$scratch/tmp" GOPATH="$scratch/gopath"
setpriv --reuid "$uid" --regid "$gid" --clear-groups --no-new-privs "$@"
