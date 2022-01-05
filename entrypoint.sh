#!/bin/sh

# This entrypoint is only used when running from GitHub Actions.

if [ -n "$1" ]; then
    exec /trojansourcedetector -config $1
else
    exec /trojansourcedetector
fi
