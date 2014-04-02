package tracker

import (
	"fmt"
	"github.com/yorirou/gotorrent/client/config"
	"github.com/yorirou/gotorrent/metainfo"
	"github.com/yorirou/gotorrent/util"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type TrackerClientCollection struct {
	clients      []*trackerClient
	infohash     string
	clientConfig *config.ClientConfig
}

func NewTrackerClientCollection(mi *metainfo.Metainfo, cc *config.ClientConfig) *TrackerClientCollection {
	tcc := new(TrackerClientCollection)
	tcc.infohash = mi.Info.Hash
	tcc.clients = tcc.createClients(mi)
	tcc.clientConfig = cc

	return tcc
}

type trackerClient struct {
	url         string
	timeout     time.Duration
	lastRequest time.Time
	collection  *TrackerClientCollection
	trackerID   string
}

func newTrackerClient(url string, tcc *TrackerClientCollection) *trackerClient {
	tc := new(trackerClient)
	tc.url = url
	tc.collection = tcc
	return tc
}

func (tc *TrackerClientCollection) createClients(mi *metainfo.Metainfo) []*trackerClient {
	clients := make([]*trackerClient, 0)

	if ann := mi.Announce; ann != "" {
		clients = append(clients, newTrackerClient(ann, tc))
	}

	for _, annurl := range mi.AnnounceList {
		clients = append(clients, newTrackerClient(annurl, tc))
	}

	return clients
}

func (tc *TrackerClientCollection) RequestPeers(downloaded, uploaded, left uint64) (seedernum uint32, leechernum uint32, peers *PeerPool) {
	var wg sync.WaitGroup

	wg.Add(len(tc.clients))
	seeders := util.NewCounter()
	leechers := util.NewCounter()
	peers = NewPeerPool()

	for _, c := range tc.clients {
		go func() {
			c.announce(tc.infohash, downloaded, uploaded, left, seeders, leechers, peers)
			wg.Done()
		}()
	}

	wg.Wait()
	seedernum = uint32(seeders.Value())
	leechernum = uint32(leechers.Value())

	return
}

func (tc *trackerClient) announce(infohash string, downloaded, uploaded, left uint64, seeders *util.Counter, leechers *util.Counter, peerPool *PeerPool) {
	if time.Since(tc.lastRequest) < tc.timeout {
		return
	}

	u, err := url.Parse(tc.url)
	if err != nil {
		log.Print(err)
		return
	}

	f := func(n uint64) string {
		return fmt.Sprintf("%d", n)
	}

	q := url.Values{}

	q.Add("info_hash", infohash)
	q.Add("peer_id", tc.collection.clientConfig.PeerID)
	q.Add("port", f(tc.collection.clientConfig.Port))
	q.Add("uploaded", f(uploaded))
	q.Add("downloaded", f(downloaded))
	q.Add("left", f(left))
	q.Add("compact", "1")
	if tc.trackerID != "" {
		u.Query().Add("trackerid", tc.trackerID)
	}

	fullurl := u.String() + "?" + q.Encode()

	resp, err := http.Get(fullurl)
	if err != nil {
		log.Print(err)
		return
	}

	b, errr := ioutil.ReadAll(resp.Body)
	if errr != nil {
		log.Print(err)
		return
	}

	r, errp := ParseResponse(b)

	if errp != nil {
		log.Print(errp)
		return
	}

	seeders.Add(uint64(r.Seeders()))
	leechers.Add(uint64(r.Leechers()))
	for _, p := range r.Peers {
		peerPool.Add(p)
	}
}
