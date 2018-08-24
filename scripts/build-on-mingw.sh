#!/usr/bin/env bash
set -e

export O3=`pwd`
export DEP_DIR=`pwd`/dep_dir
mkdir -p $DEP_DIR
export DL=/build/dls
mkdir -p $DL

echo "Deps are in: $DEP_DIR"
echo "DLs go to: $DL"
echo "O3 is: $O3"

/go-src-root/scripts/x-compile-windows.sh