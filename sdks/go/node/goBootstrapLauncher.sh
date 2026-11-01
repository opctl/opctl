#!/bin/sh
set -eu
umask 077
extract_auth() {
  auth=/run/opctl-auth
  [ "$(stat -f -c %T "$auth")" = tmpfs ] || { echo 'Go bootstrap auth requires tmpfs' >&2; exit 1; }
  head -c 1048577 > "$auth/archive"
  [ "$(wc -c < "$auth/archive")" -le 1048576 ] || { echo 'Go bootstrap private input exceeds 1 MiB' >&2; exit 1; }
  tar -tf "$auth/archive" > "$auth/list" 2>/dev/null || { echo 'Invalid Go bootstrap auth archive' >&2; exit 1; }
  while IFS= read -r name; do
    case "$name" in helper|gitconfig|.netrc) ;; *) echo 'Invalid Go bootstrap auth file' >&2; exit 1;; esac
  done < "$auth/list"
  tar -tvf "$auth/archive" 2>/dev/null | awk 'substr($0,1,1)!="-" {exit 1}' || { echo 'Go bootstrap auth files must be regular' >&2; exit 1; }
  tar --no-same-owner --no-same-permissions --keep-old-files -xf "$auth/archive" -C "$auth" 2>/dev/null || { echo 'Cannot extract Go bootstrap auth archive' >&2; exit 1; }
  rm -f "$auth/archive" "$auth/list"
  find "$auth" -type f -exec chmod 0600 {} +
  chown -R "$uid:$gid" "$auth"
}
owner=/opctl/mod
[ ! -d "$owner" ] && owner=/opctl/program
ids=$(stat -c '%u:%g' "$owner") || { echo 'Go bootstrap cache must be shared at the same host path with the container daemon' >&2; exit 1; }
uid=${ids%:*}
gid=${ids#*:}
command -v setpriv >/dev/null || { echo 'Go bootstrap requires setpriv to drop privileges' >&2; exit 1; }
scratch=$(mktemp -d /tmp/opctl-bootstrap.XXXXXX)
trap 'rm -rf "$scratch"; if [ -d /run/opctl-auth ]; then find /run/opctl-auth -mindepth 1 -delete; fi' EXIT HUP INT TERM
mkdir "$scratch/home" "$scratch/tmp" "$scratch/gopath"
chown -R "$uid:$gid" "$scratch"
if [ -d /opctl/build ]; then chown "$uid:$gid" /opctl/build; fi
export HOME="$scratch/home" TMPDIR="$scratch/tmp" GOPATH="$scratch/gopath"
if [ -d /run/opctl-auth ]; then
  extract_auth
  find /run/opctl-auth -mindepth 1 -delete
fi
setpriv --reuid "$uid" --regid "$gid" --clear-groups --no-new-privs "$@"
