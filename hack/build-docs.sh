#!/bin/sh

set -o errexit
set -o nounset

cd "$(git rev-parse --show-toplevel)" || exit 1

pip install mkdocs-material material-plausible-plugin
mkdocs build
mv site public
