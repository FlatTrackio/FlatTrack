#!/bin/bash

cd $(dirname $0)
cd $(git rev-parse --show-toplevel)

LAYOUT="${1:-}"

go install github.com/mitranim/gow@latest

if [ -z "$LAYOUT" ]; then
    echo "Options:"
    find . -name '*.kdl' | sed 's,.*/\(.*\).kdl,- \1,g'
    echo "> $0 [option]"
    exit 1
fi

zellij --layout="./hack/${LAYOUT}.kdl"
