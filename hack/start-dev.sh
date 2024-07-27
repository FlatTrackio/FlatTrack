#!/bin/bash

cd $(dirname $0)
cd $(git rev-parse --show-toplevel)

zellij --layout=./hack/dev.kdl
