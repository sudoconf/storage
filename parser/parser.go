package parser

import (
	// "errors"
	"fmt"
	"io"
	"time"

	"github.com/btlike/bencode-go"
)

// FileDict into which torrent metafile is
// parsed and stored into.
type FileDict struct {
	Length int64    `json:"length"`
	Path   []string `json:"path"`
	Md5sum string   `json:"md5sum"`
}

//InfoDict define
type InfoDict struct {
	FileDuration []int64 `json:"file-duration"`
	FileMedia    []int64 `json:"file-media"`
	// Single file
	Name   string `json:"name"`
	Length int64  `json:"length"`
	Md5sum string `json:"md5sum"`
	// Multiple files
	Files       []FileDict `json:"files"`
	PieceLength int64      `json:"piece length"`
	Pieces      string     `json:"-"`
	Private     int64      `json:"-"`
}

//MetaInfo define
type MetaInfo struct {
	Info         InfoDict   `json:"info"`
	InfoHash     string     `json:"info hash"`
	Announce     string     `json:"announce"`
	AnnounceList [][]string `json:"announce-list"`
	CreationDate int64      `json:"creation date"`
	Comment      string     `json:"comment"`
	CreatedBy    string     `json:"created by"`
	Encoding     string     `json:"encoding"`
}

//Parse Open .torrent file, un-bencode it and load them into MetaInfo struct.
func (metaInfo *MetaInfo) Parse(r io.Reader) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("bencode unmarshal panic,%v", e)
		}
	}()
	// Decode bencoded metainfo file.
	//丢弃piece部分，节省90%以上流量
	e := bencode.Unmarshal(r, metaInfo)
	if e != nil && e.Error() != "ignore piece" {
		return e
	}
	// metaInfo.DumpTorrentMetaInfo()
	return err
}

//DumpTorrentMetaInfo Print torrent meta info struct data.
func (metaInfo *MetaInfo) DumpTorrentMetaInfo() {
	fmt.Println("Announce:", metaInfo.Announce)
	fmt.Println("Announce List:")
	for _, anncListEntry := range metaInfo.AnnounceList {
		for _, elem := range anncListEntry {
			fmt.Println("    ", elem)
		}
	}
	strCreationDate := time.Unix(metaInfo.CreationDate, 0)
	fmt.Println("Creation Date:", strCreationDate)
	fmt.Println("Comment:", metaInfo.Comment)
	fmt.Println("Created By:", metaInfo.CreatedBy)
	fmt.Println("Encoding:", metaInfo.Encoding)
	fmt.Printf("InfoHash: %X\n", metaInfo.InfoHash)
	fmt.Println("Info:")
	fmt.Println("    Piece Length:", metaInfo.Info.PieceLength)
	piecesList := metaInfo.getPiecesList()
	fmt.Printf("    Pieces:%X -- %X\n", len(piecesList), len(metaInfo.Info.Pieces)/20)
	fmt.Println("    File Duration:", metaInfo.Info.FileDuration)
	fmt.Println("    File Media:", metaInfo.Info.FileMedia)
	fmt.Println("    Private:", metaInfo.Info.Private)
	fmt.Println("    Name:", metaInfo.Info.Name)
	fmt.Println("    Length:", metaInfo.Info.Length)
	fmt.Println("    Md5sum:", metaInfo.Info.Md5sum)
	fmt.Println("    Files:")
	for _, fileDict := range metaInfo.Info.Files {
		fmt.Println("        Length:", fileDict.Length)
		fmt.Println("        Path:", fileDict.Path)
		fmt.Println("        Md5sum:", fileDict.Md5sum)
	}
}

// Splits pieces string into an array of 20 byte SHA1 hashes.
func (metaInfo *MetaInfo) getPiecesList() []string {
	var piecesList []string
	piecesLen := len(metaInfo.Info.Pieces)
	for i, j := 0, 0; i < piecesLen; i, j = i+20, j+1 {
		piecesList = append(piecesList, metaInfo.Info.Pieces[i:i+19])
	}
	return piecesList
}
