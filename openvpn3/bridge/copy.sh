#!/bin/bash
#Copies header ant static library from given dir
#Usage:
# ./copy.sh <from dir> <to dir>
cp -f $1/adapter/library.h $2/.
cp -f $1/adapter/libopenvpn3.a $2/.