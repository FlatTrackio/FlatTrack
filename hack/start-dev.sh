#!/bin/bash

cd $(dirname $0)
cd $(git rev-parse --show-toplevel)

go install github.com/mitranim/gow@latest

zellij --layout=./hack/dev.kdl
