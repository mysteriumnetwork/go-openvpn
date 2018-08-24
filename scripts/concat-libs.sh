#!/usr/bin/env bash
set -e
#Usage ./concat-libs.sh outputlib.a <list of o files>
outputLib=$1
shift

mkdir -p objs
pushd objs
    for lib in ${MORE_LIBS[@]}; do
        echo "Extracting $lib"
        $AR_CMD x $lib
    done
popd

$AR_CMD rc $outputLib \
    objs/*.o \
    ${@}
rm -rf objs
