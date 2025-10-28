#!/bin/bash

cd $(dirname $0)
cd $(git rev-parse --show-toplevel)

LAYOUT="${1:-dev}"

go install github.com/mitranim/gow@latest

zellij --layout="./hack/${LAYOUT}.kdl"
