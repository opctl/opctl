apk add --update github-cli

gh api --verbose -H "Accept: application/vnd.github+json" -H "X-GitHub-Api-Version: 2022-11-28" /orgs/opctl/packages?package_type=container
