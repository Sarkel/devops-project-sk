#!/bin/bash
set -e
set -o pipefail

if [ -z "$GH_TOKEN" ] || [ -z "$PACKAGE_NAME" ] || [ -z "$OWNER" ]; then
  echo "Error: Required env vars (GH_TOKEN, PACKAGE_NAME, OWNER) are missing."
  exit 1
fi

echo "--- Processing package: $PACKAGE_NAME (Owner: $OWNER) ---"

if ! gh api "/users/$OWNER/packages/container/$PACKAGE_NAME/versions" --silent > /dev/null 2>&1; then
  echo "Warning: Package '$PACKAGE_NAME' not found or has no versions. Skipping."
  exit 0
fi

IDS_TO_DELETE=$(gh api \
  -H "Accept: application/vnd.github+json" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  "/users/$OWNER/packages/container/$PACKAGE_NAME/versions" \
  --jq '[.[] | select(.metadata.container.tags | index("latest") | not)] | .[2:] | .[].id')

if [ -z "$IDS_TO_DELETE" ]; then
  echo "No versions to delete for $PACKAGE_NAME (keeping top 2 + latest)."
  exit 0
fi

for ID in $IDS_TO_DELETE; do
  echo "Deleting version ID $ID..."
  gh api \
    --method DELETE \
    -H "Accept: application/vnd.github+json" \
    -H "X-GitHub-Api-Version: 2022-11-28" \
    "/users/$OWNER/packages/container/$PACKAGE_NAME/versions/$ID" || echo "Failed to delete $ID (continuing)"
done

echo "--- Finished $PACKAGE_NAME ---"