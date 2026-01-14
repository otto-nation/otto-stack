#!/bin/bash
set -euo pipefail

echo "🪝 Setting up Git hooks for otto-stack..."

# Create hooks directory if it doesn't exist
mkdir -p .git/hooks

# Pre-commit hook
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
set -e
if [ "$NO_VERIFY" ]; then
    echo 'pre-commit hook skipped' 1>&2
    exit 0
fi
task pre-commit
EOF

# Pre-push hook
cat > .git/hooks/pre-push << 'EOF'
#!/bin/bash
set -e

# Ensure task is in PATH for GitKraken
export PATH="/usr/local/bin:/opt/homebrew/bin:$HOME/.local/bin:$HOME/go/bin:$PATH"

if [ "$NO_VERIFY" ]; then
    echo 'pre-push hook skipped' 1>&2
    exit 0
fi
task pre-push
EOF

# Make hooks executable
chmod +x .git/hooks/pre-commit .git/hooks/pre-push

echo "✅ Git hooks installed successfully!"
echo ""
echo "Available commands:"
echo "  task pre-commit  - Run pre-commit checks manually"
echo "  task pre-push    - Run pre-push checks manually"
echo ""
echo "To skip hooks temporarily:"
echo "  git commit --no-verify"
echo "  git push --no-verify"
