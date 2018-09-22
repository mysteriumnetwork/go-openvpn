# go-openvpn

Go gettable library for wrapping openvpn functionality in go way.
There are two main options for openvpn integration:
1. As external process - go-openvpn provides external process controls (start/stop), handles management interface, can work both
as client and a service. External openvpn exe IS NOT provided (tested with openvpn 2.4.x release)
2. As built-in library - openvpn wraps c++ crosscompiled libary for all major oses (darwin,linux,win,ios and android), but has a
limitation - can only work as client only.


## Development environment

* **Step 1.** Get project dependencies
```bash
go get github.com/karalabe/xgo
go get golang.org/x/mobile/cmd/gomobile
```

* **Step 2.** Build example (Desktop)
```bash
go run examples/desktop/main.go examples/profile.ovpn
```

* **Step 3.** Build example (iOS)
```bash
./gomobile_example_ios.sh
```

## Run tests
```
make test
# We recommend running tests on frozen Linux container
scripts/xgo_run.sh make test
```

## Build
* **Step 1.** Sanity check
```bash
./check-all.sh
```

* **Step 2.** Build mobile libraries
```bash
./gomobile_ios.sh -o build/Openvpn3.framework
```