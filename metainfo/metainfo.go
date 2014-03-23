package metainfo

import (
	"github.com/yorirou/gotorrent/bencode"
)

type Metainfo struct {
	Info Info
	Announce string
	AnnounceList []string
	CreationDate uint32
	Comment string
	CreatedBy string
	Encoding string
	Name string
	Length uint64
	MD5Sum string
	Files []File
}

type Info struct {
	PieceLength uint64
	Pieces string
	Private int8
}

type File struct {
	Length uint64
	MD5Sum string
	Path []string
}

func NewMetainfo(b []byte) (*Metainfo, error) {
	mi := new(Metainfo)
	if err := bencode.Unmarshal(b, mi); err != nil {
		return nil, err
	}

	return mi, nil
}
