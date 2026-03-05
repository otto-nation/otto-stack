#!/bin/bash
set -e

echo "Checking markdown links..."

export NODE_NO_WARNINGS=1

MD_FILES=()
while IFS= read -r f; do
    MD_FILES+=("$f")
done < <(find content -name '*.md' | sort)
MD_FILES+=("README.md")

npx markdown-link-check --config .mlc_config.json "${MD_FILES[@]}" 2>/dev/null

echo "✅ All links are valid"
