package openvpn3

import "errors"

type TunnelSetup interface {
	NewBuilder() bool
	SetLayer(layer int) bool
	SetRemoteAddress(ipAddress string, ipv6 bool) bool
	AddAddress(address string, prefixLength int, gateway string, ipv6 bool, net30 bool) bool
	SetRouteMetricDefault(metric int) bool
	RerouteGw(ipv4 bool, ipv6 bool, flags int) bool
	AddRoute(address string, prefixLength int, metric int, ipv6 bool) bool
	ExcludeRoute(address string, prefixLength int, metric int, ipv6 bool) bool
	AddDnsServer(address string, ipv6 bool) bool
	AddSearchDomain(domain string) bool
	SetMtu(mtu int) bool
	SetSessionName(name string) bool
	AddProxyBypass(bypassHost string) bool
	SetProxyAutoConfigUrl(url string) bool
	SetProxyHttp(host string, port int) bool
	SetProxyHttps(host string, port int) bool
	AddWinsServer(address string) bool
	SetBlockIpv6(ipv6Block bool) bool
	SetAdapterDomainSuffix(name string) bool
	Establish() (int, error)
	Persist() bool
	EstablishLite()
	Teardown(disconnect bool)
}

//NoOpTunnelSetup by default returns false everywhere - to indicate that it is not supposed to be called and actual
//tunnel setup will not succeed
type NoOpTunnelSetup struct {
}

func (setup *NoOpTunnelSetup) NewBuilder() bool {
	return false
}

func (setup *NoOpTunnelSetup) SetLayer(layer int) bool {
	return false
}

func (setup *NoOpTunnelSetup) SetRemoteAddress(ipAddress string, ipv6 bool) bool {
	return false
}

func (setup *NoOpTunnelSetup) AddAddress(address string, prefixLength int, gateway string, ipv6 bool, net30 bool) bool {
	return false
}

func (setup *NoOpTunnelSetup) SetRouteMetricDefault(metric int) bool {
	return false
}

func (setup *NoOpTunnelSetup) RerouteGw(ipv4 bool, ipv6 bool, flags int) bool {
	return false
}

func (setup *NoOpTunnelSetup) AddRoute(address string, prefixLength int, metric int, ipv6 bool) bool {
	return false
}

func (setup *NoOpTunnelSetup) ExcludeRoute(address string, prefixLength int, metric int, ipv6 bool) bool {
	return false
}

func (setup *NoOpTunnelSetup) AddDnsServer(address string, ipv6 bool) bool {
	return false
}

func (setup *NoOpTunnelSetup) AddSearchDomain(domain string) bool {
	return false
}

func (setup *NoOpTunnelSetup) SetMtu(mtu int) bool {
	return false
}

func (setup *NoOpTunnelSetup) SetSessionName(name string) bool {
	return false
}

func (setup *NoOpTunnelSetup) AddProxyBypass(bypassHost string) bool {
	return false
}

func (setup *NoOpTunnelSetup) SetProxyAutoConfigUrl(url string) bool {
	return false
}

func (setup *NoOpTunnelSetup) SetProxyHttp(host string, port int) bool {
	return false
}

func (setup *NoOpTunnelSetup) SetProxyHttps(host string, port int) bool {
	return false
}

func (setup *NoOpTunnelSetup) AddWinsServer(address string) bool {
	return false
}

func (setup *NoOpTunnelSetup) SetBlockIpv6(ipv6Block bool) bool {
	return false
}

func (setup *NoOpTunnelSetup) SetAdapterDomainSuffix(name string) bool {
	return false
}

func (setup *NoOpTunnelSetup) Establish() (int, error) {
	return 0, errors.New("noop operation")
}

func (setup *NoOpTunnelSetup) Persist() bool {
	return false
}

func (setup *NoOpTunnelSetup) EstablishLite() {

}

func (setup *NoOpTunnelSetup) Teardown(disconnect bool) {

}

var _ TunnelSetup = &NoOpTunnelSetup{}
