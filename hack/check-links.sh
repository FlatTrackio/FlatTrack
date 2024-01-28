#!/bin/sh

lychee --exclude-path '.DS_Store' \
    ./README.md \
    ./mkdocs.yml \
    ./docs/**
