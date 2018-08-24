#!/usr/bin/env bash
set -e
rm -rf openvpn3/bridge/*.a openvpn3/bridge/*.h
docker run -it --rm -v `pwd`/build:/build -w /build -v `pwd`:/go-src-root --entrypoint /go-src-root/scripts/build-on-xgo.sh karalabe/xgo-base

#We cross compile for windows on latest ubuntu because xgo-base has older ubuntu and older mingw with windows headers which are missing
docker build -t mingw-crosscompile -f scripts/Dockerfile-mingw-ubuntu . && \
docker run -it --rm -v `pwd`/build:/build -w /build -v `pwd`:/go-src-root --entrypoint /go-src-root/scripts/build-on-mingw.sh mingw-crosscompile
