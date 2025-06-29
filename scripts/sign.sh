#!/bin/sh
set -e
if [ -z "$SIGN_KEY" ]; then
  echo "SIGN_KEY not set" >&2
  exit 1
fi
for f in "$@"; do
  gpg --batch --yes --armor --local-user "$SIGN_KEY" --output "$f.asc" --detach-sign "$f"
done
