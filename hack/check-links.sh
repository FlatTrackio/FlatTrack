#!/bin/sh

lychee --exclude-path '.DS_Store' \
    --exclude 'www\.gnu\.org' \
    ./README.md \
    ./mkdocs.yml \
    ./docs/**
