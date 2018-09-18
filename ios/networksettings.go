package ios

import (
	"ObjC/NetworkExtension"
	"ObjC/NetworkExtension/NEPacketTunnelNetworkSettings"
	"errors"
)

type NetworkSettingsCollector struct {
	remoteAddress string
}

func NewNetworkSettingsCollector() *NetworkSettingsCollector {
	return &NetworkSettingsCollector{}
}

func (collector *NetworkSettingsCollector) GetIosNetworkSettings() (NetworkExtension.NEPacketTunnelNetworkSettings, error) {
	if len(collector.remoteAddress) == 0 {
		return nil, errors.New("remote address is not set")
	}

	tunnelSettings := NEPacketTunnelNetworkSettings.NewWithTunnelRemoteAddress(collector.remoteAddress)
	return tunnelSettings, nil
}
