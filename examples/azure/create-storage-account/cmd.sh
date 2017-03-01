#!/usr/bin/env bash

echo "logging in to azure"
azure login -u "$(username)" -p "$(password)"

echo "setting default subscription"
azure account set "$(subscriptionId)"

echo "switching to ARM (azure resource manager) mode"
azure config mode arm

echo "checking for existing storage account"
if azure storage account show --resource-group "$(resourceGroupName)" "$(storageAccountName)" 1> /dev/null
then
  echo "found existing storage account"
else
  echo "creating storage account"
  azure storage account create \
    --kind "$(storageAccountKind)" \
    --location "$(location)" \
    --resource-group "$(resourceGroupName)" \
    --sku-name "$(storageAccountSkuName)" \
    "$(storageAccountName)"
fi

echo "fetching storage access keys"
storageAccessKeys=$(azure storage account keys list \
  --resource-group "$(resourceGroupName)" --json "$(storageAccountName)")

echo "storageAccessKey=$(echo $storageAccessKeys | jq -r '.[0].value')"
