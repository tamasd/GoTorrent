package config

import "github.com/yorirou/gotorrent/util"

type ClientConfig struct {
	PeerID string
	Port   uint64
}

func NewClientConfig() *ClientConfig {
	cc := new(ClientConfig)
	cc.PeerID = util.GeneratePeerID()
	return cc
}
