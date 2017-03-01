#!/usr/bin/env bash

echo "logging in to azure"
azure login -u "$USERNAME" -p "$PASSWORD"

echo "setting default subscription"
azure account set "$SUBSCRIPTION_ID"

echo "switching to ARM (azure resource manager) mode"
azure config mode arm

echo "creating resource group"
azure group create --name "$RESOURCE_GROUP_NAME" --location "$LOCATION"
