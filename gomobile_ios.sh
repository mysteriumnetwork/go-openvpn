#!/usr/bin/env bash
#Only works with patched gomobile version
CGO_LDFLAGS_ALLOW="-fobjc-arc" \
gomobile bind -work -target=ios/arm64 $@ -iosversion=10.3 -o build/IosOpenvpn3.framework github.com/mysteriumnetwork/go-openvpn/ios