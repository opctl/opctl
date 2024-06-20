apk add --update curl jq

curl --verbose \
  --user "${GITHUB_ACTOR}:${GITHUB_TOKEN}" \
  -H "Accept: application/vnd.github+json" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  -X GET https://api.github.com/orgs/opctl/packages?package_type=container | jq
