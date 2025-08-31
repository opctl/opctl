go build -tags containers_image_openpgp -buildvcs=false -o opctl-darwin-arm64 ./
sudo codesign --entitlements ./.opspec/compile/entitlements.plist --force -s - ./opctl-darwin-arm64

#opctl node delete
sudo cp -f opctl-darwin-arm64 /usr/local/bin/opctl

opctl node create