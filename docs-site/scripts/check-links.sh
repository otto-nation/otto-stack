#!/bin/bash
set -e

echo "Checking markdown links..."

# Suppress Node.js deprecation warnings and check all markdown files with config
export NODE_NO_WARNINGS=1
LINK_CHECK_OUTPUT=$(find content -name '*.md' -exec npx markdown-link-check --config .mlc_config.json {} \; 2>/dev/null && npx markdown-link-check --config .mlc_config.json README.md 2>/dev/null)
echo "$LINK_CHECK_OUTPUT"

# Check if there were any broken links (marked with ✖)
if echo "$LINK_CHECK_OUTPUT" | grep -q "\[✖\]"; then
    echo ""
    echo "❌ Link check failed: Broken links found!"
    exit 1
fi

echo "✅ All links are valid"
