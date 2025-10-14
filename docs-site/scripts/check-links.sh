#!/bin/bash
set -e

echo "Checking markdown links..."

# Suppress Node.js deprecation warnings and check all markdown files
export NODE_NO_WARNINGS=1
LINK_CHECK_OUTPUT=$(find content -name '*.md' -exec markdown-link-check {} \; 2>/dev/null && markdown-link-check README.md 2>/dev/null)
echo "$LINK_CHECK_OUTPUT"

# Check if there were any broken links (marked with ✖)
if echo "$LINK_CHECK_OUTPUT" | grep -q "\[✖\]"; then
    echo ""
    echo "❌ Link check failed: Broken links found!"
    exit 1
fi

echo "✅ All links are valid"
