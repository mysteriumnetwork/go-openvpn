#!/usr/bin/env bash

CGO_LDFLAGS_ALLOW="-fobjc-arc" \
gomobile bind -target=ios/arm64 $@ -iosversion=10.3 -v github.com/mysteriumnetwork/go-openvpn/openvpn3