package ios

import (
	"ObjC/NetworkExtension"
	"ObjC/NetworkExtension/NEPacket"
)

type IOSPacketFlowTunnel struct {
}

func NewIOSPacketFlowTunnel(flow NetworkExtension.NEPacketTunnelFlow) *IOSPacketFlowTunnel {
	data := []byte{0x01}
	_ = NEPacket.NewWithData(data, 1)
	//_ := NSArray.ArrayWithObject(packet)
	//_ = NSArray.Array()
	//_ = flow.WritePacketObjects(arr)
	return &IOSPacketFlowTunnel{}
}
