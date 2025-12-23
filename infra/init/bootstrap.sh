#!/bin/bash
set -e

ENV_FILE=".env"

echo "Starting infrastructure bootstrapping (Robust Mode)..."

if [ -f "$ENV_FILE" ]; then
    echo "Loading configuration from $ENV_FILE..."
    set -a
    source "$ENV_FILE"
    set +a
else
    echo "ERROR: File $ENV_FILE not found!"
    exit 1
fi

echo "Detecting active subscription..."
SUBSCRIPTION_ID=$(az account show --query id --output tsv)

if [ -z "$SUBSCRIPTION_ID" ]; then
    echo "ERROR: No active subscription found. Run 'az login' first."
    exit 1
fi
echo "Using Subscription ID: $SUBSCRIPTION_ID"

echo "Registering Microsoft.Storage provider..."
az provider register --namespace Microsoft.Storage --subscription "$SUBSCRIPTION_ID"
# Give it a moment to propagate
sleep 5

echo "Creating Resource Group: $RESOURCE_GROUP_NAME..."
az group create \
    --name "$RESOURCE_GROUP_NAME" \
    --location "$LOCATION" \
    --subscription "$SUBSCRIPTION_ID"

echo "Creating Storage Account: $STORAGE_ACCOUNT_NAME..."
az storage account create \
    --resource-group "$RESOURCE_GROUP_NAME" \
    --name "$STORAGE_ACCOUNT_NAME" \
    --sku Standard_LRS \
    --encryption-services blob \
    --location "$LOCATION" \
    --subscription "$SUBSCRIPTION_ID"

echo "Creating Blob Container: $CONTAINER_NAME..."
# We need to grab the storage key first to avoid login/auth issues during container creation
STORAGE_KEY=$(az storage account keys list --resource-group "$RESOURCE_GROUP_NAME" --account-name "$STORAGE_ACCOUNT_NAME" --subscription "$SUBSCRIPTION_ID" --query "[0].value" -o tsv)

az storage container create \
    --name "$CONTAINER_NAME" \
    --account-name "$STORAGE_ACCOUNT_NAME" \
    --account-key "$STORAGE_KEY" \
    --subscription "$SUBSCRIPTION_ID"

echo "Creating Service Principal (RBAC)..."
# We use the captured SUBSCRIPTION_ID
SP_JSON=$(az ad sp create-for-rbac \
    --name "sp-devops-project-github-actions" \
    --role="Contributor" \
    --scopes="/subscriptions/$SUBSCRIPTION_ID" \
    --output json)

APP_ID=$(echo $SP_JSON | jq -r .appId)
PASSWORD=$(echo $SP_JSON | jq -r .password)
TENANT_ID=$(echo $SP_JSON | jq -r .tenant)

echo ""
echo "SUCCESS! Environment ready."
echo "========================================="
echo "Copy the following values to GitHub Repository Secrets:"
echo "-----------------------------------------"
echo "AZURE_SUBSCRIPTION_ID : $SUBSCRIPTION_ID"
echo "AZURE_TENANT_ID       : $TENANT_ID"
echo "AZURE_CLIENT_ID       : $APP_ID"
echo "AZURE_CLIENT_SECRET   : $PASSWORD"
echo "-----------------------------------------"
echo "Update infra/main.tf (backend block):"
echo "resource_group_name  = \"$RESOURCE_GROUP_NAME\""
echo "storage_account_name = \"$STORAGE_ACCOUNT_NAME\""
echo "container_name       = \"$CONTAINER_NAME\""
echo "key                  = \"\${env}.terraform.tfstate\""
echo "========================================="