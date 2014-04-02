package tracker

import (
	"encoding/binary"
	"fmt"
	"github.com/yorirou/gotorrent/bencode"
	"log"
)

func ParseResponse(resp []byte) (*Response, error) {
	compact := new(CompactResponse)
	if err := bencode.Unmarshal(resp, compact); err != nil {
		log.Print(err)
		dictcompact := new(Response)
		if errd := bencode.Unmarshal(resp, dictcompact); err != nil {
			return nil, errd
		}
		return dictcompact, nil
	}

	return compact.Convert(), nil
}

type ResponseBase struct {
	FailureReason  string
	WarningMessage string
	Interval       uint64
	MinInterval    uint64
	TrackerID      string
	Complete       uint32
	Incomplete     uint32
}

func (rb *ResponseBase) GetInterval(min bool) uint64 {
	if min {
		return rb.MinInterval
	}

	return rb.Interval
}

func (rb *ResponseBase) Seeders() uint32 {
	return rb.Complete
}

func (rb *ResponseBase) Leechers() uint32 {
	return rb.Incomplete
}

type CompactResponse struct {
	ResponseBase
	Peers []byte
	Peers6 []byte
}

func (cr *CompactResponse) Convert() *Response {
	dr := new(Response)
	dr.ResponseBase = cr.ResponseBase

	peerdata := []byte(cr.Peers)

	if len(peerdata)%6 != 0 {
		log.Print("peer data is not divisible with 6")
	}

	peernum := len(peerdata) / 6
	dr.Peers = make([]*Peer, peernum)
	for i := 0; i < peernum; i++ {
		p := new(Peer)

		p.PeerID = ""
		p.IP = fmt.Sprintf("%d.%d.%d.%d", peerdata[i*6+0], peerdata[i*6+1], peerdata[i*6+2], peerdata[i*6+3])
		p.Port = binary.BigEndian.Uint16(peerdata[i*6+4 : i*6+6])

		dr.Peers[i] = p
	}

	return dr
}

type Response struct {
	ResponseBase
	Peers []*Peer
}
