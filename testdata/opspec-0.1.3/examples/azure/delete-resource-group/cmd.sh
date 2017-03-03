#!/usr/bin/env bash

echo "logging in to azure"
azure login -u "$USERNAME" -p "$PASSWORD"

echo "setting default subscription"
azure account set "$SUBSCRIPTION_ID"

echo "switching to ARM (azure resource manager) mode"
azure config mode arm

echo "checking for existing resource group"
if azure group show "$RESOURCE_GROUP_NAME" 1> /dev/null
then
  echo "deleting resource group"
  azure group delete --name "$RESOURCE_GROUP_NAME" --quiet
else
  echo "existing resource group not found"
fi
