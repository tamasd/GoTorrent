package main

import (
	"flag"
	"fmt"
	"github.com/yorirou/gotorrent/client/config"
	"github.com/yorirou/gotorrent/metainfo"
	"github.com/yorirou/gotorrent/torrent"
	"github.com/yorirou/gotorrent/tracker"
	"github.com/yorirou/gotorrent/util"
	"io/ioutil"
	"log"
	"os"
	"runtime"
)

var action = flag.String("action", "", "info, announce, download")

func main() {
	flag.Parse()

	// Set up the logger
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Configure the scheduler
	runtime.GOMAXPROCS(runtime.NumCPU())

	actions := map[string]func([]byte){
		"info":     info,
		"announce": announce,
		"download": download,
	}

	args := flag.Args()

	if len(args) != 1 {
		log.Fatal("1 argument is allowed, which is the .torrent file")
	}

	file, ferr := os.Open(args[0])
	if ferr != nil {
		log.Fatal(ferr)
	}

	fc, ioerr := ioutil.ReadAll(file)
	if ioerr != nil {
		log.Fatal(ioerr)
	}

	callback, ok := actions[*action]
	if !ok {
		log.Fatal("action must be info or announce or download")
	}
	callback(fc)
}

func announce(torrentfile []byte) {
	mi, err := metainfo.NewMetainfo(torrentfile)
	if err != nil {
		log.Fatal(err)
	}

	cfg := config.NewClientConfig()
	cfg.PeerID = util.GeneratePeerID()
	cfg.Port = 7000

	t := torrent.NewTorrent(mi, cfg)

	tcc := tracker.NewTrackerClientCollection(mi, cfg)
	seeders, leechers, peers := tcc.RequestPeers(t.Downloaded(), t.Uploaded(), t.Left())

	fmt.Printf("Seeders: %d\nLeechers: %d\nPeers: %v\n", seeders, leechers, peers)
}

func info(torrentfile []byte) {
	mi, err := metainfo.NewMetainfo(torrentfile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(mi)
}

func download(torrentfile []byte) {

}
