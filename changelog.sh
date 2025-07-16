#!/bin/bash

set -e

startTag=$(git tag --list --sort=version:refname | tail -n 1)
lastTag=$(git tag --list --sort=version:refname | tail -n 2 | head -n 1)
#
# startTag=HEAD
# lastTag=cf40913

if [ -z "$startTag" ]; then
    echo "No tags found. Please create a tag first."
    exit 1
fi

if [ -z "$lastTag" ]; then
    echo "Only one tag found: $startTag. No changes to show."
    exit 0
fi

date=$(git log -1 --format=%cd --date=short "$startTag")

echo "## [${startTag}](https://github.com/chriserin/seq/compare/${startTag}...${lastTag}) ($date)"
echo ""

features=$(git log --oneline "${startTag}"..."${lastTag}" | awk ' $2 ~ /^feat/ {print}')

if [ -n "$features" ]; then
    echo "### Features"
    echo ""
    while IFS= read -r line; do
        commit_hash=$(echo "$line" | awk '{print $1}')
        commit_message=$(echo "$line" | cut -d' ' -f3-)
        echo "* $commit_message [${commit_hash}](https://github.com/chriserin/seq/commit/${commit_hash}) "
    done <<<"$features"
else
    echo "### Features"
    echo ""
    echo "No new features."
    echo ""
fi

fixes=$(git log --oneline "${startTag}"..."${lastTag}" | awk ' $2 ~ /^fix/ {print}')

if [ -n "$fixes" ]; then
    echo ""
    echo "### Fixes"
    echo ""
    while IFS= read -r line; do
        commit_hash=$(echo "$line" | awk '{print $1}')
        commit_message=$(echo "$line" | cut -d' ' -f3-)
        echo "* $commit_message [${commit_hash}](https://github.com/chriserin/seq/commit/${commit_hash}) "
    done <<<"$fixes"
else
    echo "### Fixes"
    echo ""
    echo "No bug fixes."
    echo ""
fi
