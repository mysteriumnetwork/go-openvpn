# go-openvpn

Go gettable library for wrapping openvpn functionality in go way.
There are two main options for openvpn integration:
1. As external process - go-openvpn provides external process controls (start/stop), handles management interface, can work both
as client and a service. External openvpn exe IS NOT provided (tested with openvpn 2.4.x release)
2. As built-in library - openvpn wraps c++ crosscompiled libary for all major oses (darwin,linux,win,ios and android), but has a
limitation - can only work as client only.
