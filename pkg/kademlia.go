package godemlia

import "strings"

type NodeInfo struct {
	ID   []byte `json:"id,omitempty"`
	IP   string `json:"ip,omitempty"`
	Port int    `json:"port,omitempty"`
}

func (b *NodeInfo) Equal(other NodeInfo) bool {
	return strings.EqualFold(b.IP, other.IP) &&
		b.Port == other.Port
}
