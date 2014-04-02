package torrent

import (
	"github.com/yorirou/gotorrent/client/config"
	"github.com/yorirou/gotorrent/metainfo"
	"github.com/yorirou/gotorrent/tracker"
)

type Torrent struct {
	metainfo   *metainfo.Metainfo
	uploaded   uint64
	downloaded uint64
	seeders    uint32
	leechers   uint32
	peers      *tracker.PeerPool
	trackers   *tracker.TrackerClientCollection
}

func NewTorrent(mi *metainfo.Metainfo, cc *config.ClientConfig) *Torrent {
	t := new(Torrent)
	t.metainfo = mi
	t.trackers = tracker.NewTrackerClientCollection(mi, cc)
	return t
}

func (t *Torrent) RequestPeers() {
	t.seeders, t.leechers, t.peers = t.trackers.RequestPeers(t.downloaded, t.uploaded, t.Left())
}

func (t *Torrent) GetMetaInfo() *metainfo.Metainfo {
	return t.metainfo
}

func (t *Torrent) Left() uint64 {
	return t.metainfo.Length - t.downloaded
}

func (t *Torrent) ResetUploaded() {
	t.uploaded = 0
}

func (t *Torrent) ResetDownloaded() {
	t.downloaded = 0
}

func (t *Torrent) AddToUploaded(bytes uint64) {
	t.uploaded += bytes
}

func (t *Torrent) AddToDownloaded(bytes uint64) {
	t.downloaded += bytes
}

func (t *Torrent) Uploaded() uint64 {
	return t.uploaded
}

func (t *Torrent) Downloaded() uint64 {
	return t.downloaded
}
