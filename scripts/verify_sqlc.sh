#!/bin/bash
# Exits immediately if a command exits with a non-zero status
set -e

echo "=== STARTING SQLC VERIFICATION ==="

echo " > Stashing current working directory state (including script chmod)..."
git stash push -u -m "CI_STATE_BEFORE_SQLC" > /dev/null

echo " > Regenerating SQLC code..."
make sqlc > /dev/null

CHANGES=$(git status --porcelain)

if [ -n "$CHANGES" ]; then
    # --- FAILURE PATH ---

    echo "::error::DRIFT DETECTED! The generated Go code is out of sync with SQL queries."
    echo "--------------------------------------------------"
    echo "The following files were modified/created by SQLC:"
    echo "$CHANGES"
    echo "--------------------------------------------------"
    echo "ACTION REQUIRED: Run 'make sqlc' locally and commit the changes."

    # Restore the original state of the working directory before exiting
    # '|| true' makes the step robust against minor warnings if the stash was empty
    git stash pop --index > /dev/null || true
    exit 1
else
    # --- SUCCESS PATH ---

    echo "Success: SQLC code is clean and in sync with the repository."

    # Restore the original working directory state (this restores the executable bit on the script)
    git stash pop --index > /dev/null || true
fi