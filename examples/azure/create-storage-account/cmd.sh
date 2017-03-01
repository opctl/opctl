#!/usr/bin/env bash

echo "logging in to azure"
azure login -u "$USERNAME" -p "$PASSWORD"

echo "setting default subscription"
azure account set "$SUBSCRIPTION_ID"

echo "switching to ARM (azure resource manager) mode"
azure config mode arm

echo "checking for existing storage account"
if azure storage account show --resource-group "$RESOURCE_GROUP_NAME" "$STORAGE_ACCOUNT_NAME" 1> /dev/null
then
  echo "found existing storage account"
else
  echo "creating storage account"
  azure storage account create \
    --kind "$STORAGE_ACCOUNT_KIND" \
    --location "$LOCATION" \
    --resource-group "$RESOURCE_GROUP_NAME" \
    --sku-name "$STORAGE_ACCOUNT_SKU_NAME" \
    "$STORAGE_ACCOUNT_NAME"
fi

echo "fetching storage access keys"
storageAccessKeys=$(azure storage account keys list \
  --resource-group "$RESOURCE_GROUP_NAME" --json "$STORAGE_ACCOUNT_NAME")

echo "storageAccessKey=$(echo $storageAccessKeys | jq -r '.[0].value')"
