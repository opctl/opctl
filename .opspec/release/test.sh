apk add --update curl jq

curl --verbose \
  -H "Accept: application/vnd.github+json" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  -H "Authorization: Bearer ${GH_TOKEN}"\
  https://api.github.com/orgs/opctl/packages?package_type=container | jq
