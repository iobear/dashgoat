#!/bin/bash

echo
echo "-- link bin --"

SOURCE_DIR="build/"

# Determine the operating system
OS=$(uname -s)

# Determine the architecture
ARCH=$(uname -m)


echo "Found: $OS $ARCH"
DGBIN=""

if [[ "$OS" == "Darwin" ]]; then
    if [[ "$ARCH" == "x86_64" ]]; then
        DGBIN="dashgoat-mac"
    else
        DGBIN="dashgoat-mac-arm"
    fi
fi

if [[ "$OS" == "Linux" ]]; then
    if [[ "$ARCH" == "x86_64" ]]; then
        DGBIN="dashgoat"
    fi
fi

if [[ "$DGBIN" == "" ]]; then
    echo "Unsupported system: $OS $ARCH"
fi

rm -f dashgoat
ln -s $SOURCE_DIR$DGBIN dashgoat

#Debug
pwd
ls
ls build