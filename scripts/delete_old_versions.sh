#!/bin/bash
set -e
set -o pipefail

if [ -z "$GITHUB_REPOSITORY" ]; then
  echo "Error: GITHUB_REPOSITORY env var is missing. Are you running this locally?"
  exit 1
fi

if [ -z "$GH_TOKEN" ] || [ -z "$PACKAGE_NAME" ] || [ -z "$OWNER" ]; then
  echo "Error: Required env vars (GH_TOKEN, PACKAGE_NAME, OWNER) are missing."
  exit 1
fi

PACKAGE_NAME_ENCODED=${PACKAGE_NAME////%2F}

echo "--- Processing package: $PACKAGE_NAME (Encoded: $PACKAGE_NAME_ENCODED) ---"
echo "--- Context: Owner: $OWNER | Repo: $GITHUB_REPOSITORY ---"

SCOPE="users"
if gh api "/orgs/$OWNER" --silent > /dev/null 2>&1; then
  SCOPE="orgs"
fi

if ! gh api "/$SCOPE/$OWNER/packages/container/$PACKAGE_NAME_ENCODED/versions" --silent > /dev/null 2>&1; then
  echo "Warning: Package '$PACKAGE_NAME' NOT found."

  echo "Listing packages linked to repository '$GITHUB_REPOSITORY' to show valid names:"

  REPO_PACKAGES=$(gh api "/repos/$GITHUB_REPOSITORY/packages?package_type=container" --jq '.[].name' 2>/dev/null || echo "")

  if [ -z "$REPO_PACKAGES" ]; then
     echo "   (No packages found linked to this repository)"
  else
     echo "$REPO_PACKAGES"
     echo ""
     echo "💡 HINT: Update your 'cleanup.yml' matrix to match one of the names above exactly."
  fi

  exit 0
fi

IDS_TO_DELETE=$(gh api \
  -H "Accept: application/vnd.github+json" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  "/$SCOPE/$OWNER/packages/container/$PACKAGE_NAME_ENCODED/versions" \
  --jq '[.[] | select(.metadata.container.tags | index("latest") | not)] | .[2:] | .[].id')

if [ -z "$IDS_TO_DELETE" ]; then
  echo "No old versions to delete for $PACKAGE_NAME."
  exit 0
fi

for ID in $IDS_TO_DELETE; do
  echo "Deleting version ID $ID..."
  gh api \
    --method DELETE \
    -H "Accept: application/vnd.github+json" \
    -H "X-GitHub-Api-Version: 2022-11-28" \
    "/$SCOPE/$OWNER/packages/container/$PACKAGE_NAME_ENCODED/versions/$ID" || echo "❌ Failed to delete $ID (continuing)"
done

echo "--- Finished $PACKAGE_NAME ---"