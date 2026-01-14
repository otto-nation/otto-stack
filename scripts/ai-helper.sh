#!/bin/bash
set -euo pipefail

get_ai_command() {
  if [ -f .env ] && grep -q "^AI_COMMAND=" .env; then
    grep "^AI_COMMAND=" .env | cut -d'=' -f2-
  else
    return 1
  fi
}

check_ai_available() {
  AI_COMMAND=$(get_ai_command)
  if [ $? -eq 0 ] && [ -n "$AI_COMMAND" ]; then
    AI_CMD=$(echo "$AI_COMMAND" | cut -d' ' -f1)
    if command -v "$AI_CMD" >/dev/null 2>&1; then
      echo "$AI_COMMAND"
      return 0
    else
      echo "❌ AI command not found: $AI_CMD" >&2
      return 1
    fi
  else
    echo "❌ AI not configured. Please set AI_COMMAND in .env file" >&2
    echo "   Example: AI_COMMAND=kiro-cli --no-interactive" >&2
    return 1
  fi
}
