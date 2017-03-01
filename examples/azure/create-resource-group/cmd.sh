#!/usr/bin/env bash

echo "logging in to azure"
azure login -u "$(username)" -p "$(password)"

echo "setting default subscription"
azure account set "$(subscriptionId)"

echo "switching to ARM (azure resource manager) mode"
azure config mode arm

echo "creating resource group"
azure group create --name "$(resourceGroupName)" --location "$(location)"
