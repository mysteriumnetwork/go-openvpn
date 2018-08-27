#!/usr/bin/env bash
set -e
rm -rf openvpn3/bridge/*.a openvpn3/bridge/*.h
docker run -it --rm -v `pwd`/build:/build -w /build -v `pwd`:/go-src-root --entrypoint /go-src-root/scripts/build-on-xgo.sh mysteriumnetwork/xgo-1.9.2
