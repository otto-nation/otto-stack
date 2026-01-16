# Homebrew Tap

Automated deployment of otto-stack to the Homebrew tap repository.

## Overview

When a release is published, the workflow automatically:
1. Generates the Homebrew formula with checksums
2. Pushes the updated formula to `otto-nation/homebrew-tap`

Users install with:
```bash
brew install otto-nation/tap/otto-stack
```

## Setup

Add `HOMEBREW_TAP_TOKEN` to repository secrets:

1. Create PAT at https://github.com/settings/tokens/new with `repo` scope
2. Add to https://github.com/otto-nation/otto-stack/settings/secrets/actions
3. Name: `HOMEBREW_TAP_TOKEN`

## Manual Deployment

```bash
# Update formula with checksums
./scripts/homebrew.sh update -v v1.2.3

# Deploy to tap (requires token)
HOMEBREW_TAP_TOKEN=your_token ./scripts/homebrew.sh deploy -v v1.2.3

# Dry run
./scripts/homebrew.sh deploy -v v1.2.3 --dry-run
```

## Token Priority

1. `HOMEBREW_TAP_TOKEN` (preferred)
2. `GITHUB_TOKEN` (fallback, limited permissions)
