package metainfo

import (
	"encoding/base64"
	"fmt"
	"github.com/yorirou/gotorrent/bencode"
	"github.com/yorirou/gotorrent/util"
	"strings"
	"time"
)

type Metainfo struct {
	Info         Info
	Announce     string
	AnnounceList []string
	CreationDate uint32
	Comment      string
	CreatedBy    string
	Encoding     string
	Name         string
	Length       uint64
	MD5Sum       string
	Files        []File
	HTTPSeeds    []string
}

type Info struct {
	Hash        string
	PieceLength uint64
	Pieces      []byte
	Private     int8
	Length      uint64
	Name        string
}

type File struct {
	Length uint64
	MD5Sum string
	Path   []string
}

func NewMetainfo(b []byte) (*Metainfo, error) {
	mi := new(Metainfo)
	if err := bencode.UnmarshalWithRaw(b, mi, "Hash"); err != nil {
		return nil, err
	}

	mi.Info.Hash = util.Hash(mi.Info.Hash)

	return mi, nil
}

func (mi *Metainfo) String() string {
	output := ""

	output += "Info:\n\t" + strings.Replace(mi.Info.String(), "\n", "\n\t", -1) + "\n"
	output += "Announce: " + mi.Announce + "\n"
	output += "AnnounceList: \n"
	for _, ali := range mi.AnnounceList {
		output += "\t" + ali + "\n"
	}
	output += "CreationDate: " + time.Unix(int64(mi.CreationDate), 0).Format(time.RFC3339) + "\n"
	output += "Comment: " + mi.Comment + "\n"
	output += "CreatedBy: " + mi.CreatedBy + "\n"
	output += "Encoding: " + mi.Encoding + "\n"
	output += "Name: " + mi.Name + "\n"
	output += fmt.Sprintf("Length: %d\n", mi.Length)
	output += "MD5Sum: " + mi.MD5Sum + "\n"
	output += "Files: \n\t"
	for _, f := range mi.Files {
		output += strings.Replace(f.String(), "\n", "\n\t", -1)
	}
	output += "\n"
	output += "HTTPSeeds: " + strings.Join(mi.HTTPSeeds, ", ") + "\n"

	return output
}

func (i *Info) String() string {
	output := ""

	output += "Hash: " + base64.URLEncoding.EncodeToString([]byte(i.Hash)) + "\n"
	output += fmt.Sprintf("PiecesLength: %d\n", i.PieceLength)
	output += "Pieces: " + base64.URLEncoding.EncodeToString(i.Pieces) + "\n"
	output += fmt.Sprintf("Private: %d\n", i.Private)
	output += fmt.Sprintf("Length: %d\n", i.Length)
	output += "Name: " + i.Name + "\n"

	return output
}

func (f *File) String() string {
	output := ""

	output += fmt.Sprintf("Length: %d\n", f.Length)
	output += "MD5Sum: " + f.MD5Sum + "\n"
	output += "Path: " + strings.Join(f.Path, "/")

	return output
}
