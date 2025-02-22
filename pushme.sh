#!/bin/bash

# Check if a commit message is provided
if [ -z "$1" ]; then
  echo "Error: Commit message is required."
  echo "Usage: ./pushme.sh \"Your commit message\""
  exit 1
fi

# Commit message from the argument
COMMIT_MESSAGE="$1"

# Stage all changes
git add .

# Commit with the provided message
git commit -m "$COMMIT_MESSAGE"

# Get the name of the current branch
CURRENT_BRANCH=$(git symbolic-ref --short HEAD)

# Check if the current branch has an upstream
UPSTREAM=$(git rev-parse --abbrev-ref "${CURRENT_BRANCH}@{upstream}" 2>/dev/null)

# If no upstream is set, set it to the origin of the current branch
if [ -z "$UPSTREAM" ]; then
  echo "No upstream set for branch $CURRENT_BRANCH. Setting upstream to origin/$CURRENT_BRANCH."
  git push --set-upstream origin "$CURRENT_BRANCH"
else
  # If upstream is set, just push to the current branch
  git push
fi
