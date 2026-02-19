#!/bin/sh

set -e

# Configure git to use hooks from .githooks directory
git config --local core.hooksPath .githooks/

# Make hooks executable
chmod +x .githooks/pre-commit
chmod +x .githooks/post-checkout

echo "✅ Git hooks installed successfully."
echo "   - pre-commit: Runs compile check, tests, and linter"
echo "   - post-checkout: Installs dependencies and validates setup"