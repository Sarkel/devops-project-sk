#!/bin/bash
set -e

# Validate required variables
if [ -z "$ACR_NAME" ] || [ -z "$REPOSITORY" ]; then
  echo "Error: ACR_NAME and REPOSITORY env vars are required."
  exit 1
fi

echo "--- Processing ACR Repository: $REPOSITORY in $ACR_NAME ---"

if ! az acr repository show --name "$ACR_NAME" --repository "$REPOSITORY" > /dev/null 2>&1; then
  echo "Repository '$REPOSITORY' not found in registry '$ACR_NAME'. Skipping."
  exit 0
fi

DIGESTS_TO_DELETE=$(az acr repository show-manifests \
  --name "$ACR_NAME" \
  --repository "$REPOSITORY" \
  --orderby time_desc \
  --query "[].digest" \
  -o tsv \
  | tail -n +4)

if [ -z "$DIGESTS_TO_DELETE" ]; then
  echo "No old manifests to delete for $REPOSITORY (keeping top 3)."
  exit 0
fi

echo "$DIGESTS_TO_DELETE" | while read -r DIGEST; do
  if [ -n "$DIGEST" ]; then
    echo "Deleting manifest: $DIGEST"
    # --yes skips the confirmation prompt
    az acr repository delete --name "$ACR_NAME" --image "$REPOSITORY@$DIGEST" --yes
  fi
done

echo "--- Finished $REPOSITORY ---"