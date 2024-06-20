apk add --update docker

error=$(docker pull ghcr.io/opctl/opctl:${version}-dind 2>&1 1>/dev/null)

# only continue with release if tag doesn't exist
if [[ -z $error ]]; then
  echo "Opctl Image for version '${version}' already exists"
  echo -n true > /alreadyPublished
elif [[ $error =~ "denied" ]]; then
  echo "unable to read from ghcr with given credentials"
  exit 1
else
  echo "error was ${error}"
  echo "Image does not exist for version '${version}', proceeding with release..."
  echo -n false > /alreadyPublished
fi
