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
  download_status=0
  setpriv --reuid "$uid" --regid "$gid" --clear-groups --no-new-privs sh -c '
    export HOME=/run/opctl-auth GIT_CONFIG_GLOBAL=/run/opctl-auth/gitconfig GIT_CONFIG_NOSYSTEM=1 GIT_TERMINAL_PROMPT=0
    export GIT_TRACE=0 GIT_TRACE_CURL=0 GIT_CURL_VERBOSE=0 GIT_TRACE_PACKET=0 GIT_TRACE_SETUP=0 GIT_TRACE_PERFORMANCE=0 GIT_SSH_COMMAND=false
    mkfifo "$HOME/stdout.pipe" "$HOME/stderr.pipe"
    head -c 65537 < "$HOME/stdout.pipe" > "$HOME/download.stdout" & out=$!
    head -c 65537 < "$HOME/stderr.pipe" > "$HOME/download.stderr" & err=$!
    status=0
    go mod download > "$HOME/stdout.pipe" 2> "$HOME/stderr.pipe" || status=$?
    wait "$out" "$err"
    rm -f "$HOME/stdout.pipe" "$HOME/stderr.pipe"
    exit "$status"
  ' || download_status=$?
  for stream in stdout stderr; do
    [ -f "/run/opctl-auth/download.$stream" ] || { echo 'Go module download failed before diagnostic capture' >&2; exit 1; }
    [ "$(wc -c < "/run/opctl-auth/download.$stream")" -le 65536 ] || { echo 'Go module download diagnostics exceed 64 KiB' >&2; exit 1; }
  done
  if [ "$download_status" -ne 0 ]; then
    if grep -Eiq '401|403|unauthorized|authentication|terminal prompts disabled' /run/opctl-auth/download.stderr; then
      echo 'Go module download/auth failed. Repository-scoped credentials support HTTPS Git; authenticated vanity discovery on Go 1.21 needs host-scoped auth.' >&2
    else
      echo 'Go module download failed. Check auth and update go.mod/go.sum outside bootstrap; private download diagnostics are suppressed.' >&2
    fi
    exit 1
  fi
  find /run/opctl-auth -mindepth 1 -delete
elif [ -d /opctl/build ]; then
  setpriv --reuid "$uid" --regid "$gid" --clear-groups --no-new-privs go mod download
fi
if [ -d /opctl/build ]; then
  export GOPROXY=off GONOPROXY=none
fi
setpriv --reuid "$uid" --regid "$gid" --clear-groups --no-new-privs "$@"
