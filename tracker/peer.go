package tracker

import (
	"fmt"
	"sync"
)

type Peer struct {
	PeerID string
	IP     string
	Port   uint16
}

func (p *Peer) Hash() string {
	return fmt.Sprintf("tracker.Peer:%s:%d", p.IP, p.Port)
}

func (p *Peer) String() string {
	return fmt.Sprintf("PeerID: %s IP: %s Port: %d", p.PeerID, p.IP, p.Port)
}

type PeerPool struct {
	peers map[string]*Peer
	mtx   sync.RWMutex
}

func NewPeerPool() *PeerPool {
	pp := new(PeerPool)
	pp.peers = make(map[string]*Peer)

	return pp
}

func (pp *PeerPool) String() string {
	output := ""
	for _, peer := range pp.peers {
		output += peer.String() + "\n"
	}
	return output
}

func (pp *PeerPool) Add(p *Peer) {
	pp.mtx.Lock()
	defer pp.mtx.Unlock()
	pp.peers[p.Hash()] = p
}

func (pp *PeerPool) GetPeers() []*Peer {
	pp.mtx.RLock()
	defer pp.mtx.RUnlock()

	peers := make([]*Peer, len(pp.peers))
	i := 0
	for _, p := range pp.peers {
		peers[i] = p
		i++
	}

	return peers
}
