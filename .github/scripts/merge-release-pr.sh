#!/usr/bin/env bash

# Merges the release PR that release-please keeps open, cutting a release for
# whatever has landed on main since the last one. This is what gets dependency
# bumps out to consumers such as Athenaeum, which otherwise would only see a new
# version when a feature happened to be released.
#
# This must run with a PAT rather than GITHUB_TOKEN: a push made with
# GITHUB_TOKEN does not trigger workflows, so the main build would not run and
# release-please would never tag the release.

set -euo pipefail

pr_number=$(gh pr list --label "autorelease: pending" --state open --json number --jq '.[0].number // empty')

if [ -z "$pr_number" ]; then
    echo "No pending release PR, nothing to release"
    exit 0
fi

echo "Merging release PR #${pr_number}"
gh pr merge "$pr_number" --squash --delete-branch
