#!/usr/bin/env bash

if command -v <go> &> /dev/null
then
    go build -o $PREFIX/chhome main.go
else
    echo "Go is not detected. Make sure Go is in your PATH"
    exit 1
fi
