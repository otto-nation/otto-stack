#!/bin/bash
set -e

echo "Checking markdown links..."

export NODE_NO_WARNINGS=1

# Collect all markdown files into a single invocation so each URL is checked
# once rather than once per file, avoiding per-domain rate limiting.
MD_FILES=()
while IFS= read -r f; do
    MD_FILES+=("$f")
done < <(find content -name '*.md' | sort)
MD_FILES+=("README.md")

LINK_CHECK_OUTPUT=$(npx markdown-link-check --config .mlc_config.json "${MD_FILES[@]}" 2>/dev/null)
echo "$LINK_CHECK_OUTPUT"

if echo "$LINK_CHECK_OUTPUT" | grep -q "\[✖\]"; then
    echo ""
    echo "❌ Link check failed: Broken links found!"
    exit 1
fi

echo "✅ All links are valid"
