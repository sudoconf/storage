package crawl

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/btlike/database/torrent"
	"github.com/btlike/storage/parser"
	"github.com/btlike/storage/utils"
)

type torrentData struct {
	Infohash   string
	Name       string
	CreateTime time.Time
	Length     int64
	FileCount  int64

	Files []file
}

type files []file

type file struct {
	Name   string
	Length int64
}

func (a files) Len() int           { return len(a) }
func (a files) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a files) Less(i, j int) bool { return a[i].Length > a[j].Length }

//Store data in database
func Store(data parser.MetaInfo) (err error) {
	var t torrentData
	t.Infohash = data.InfoHash
	if len(t.Infohash) != 40 {
		return fmt.Errorf("store infohash len is not 40")
	}
	if data.Info.Name == "" {
		// fmt.Println("store name len is 0")
		return fmt.Errorf("store name len is 0")
	}
	t.Name = data.Info.Name
	t.CreateTime = time.Now()
	if len(data.Info.Files) == 0 {
		t.Length = data.Info.Length
		t.FileCount = 1
		t.Files = append(t.Files, file{Name: t.Name, Length: t.Length})
	} else {
		var tmpFiles files
		if len(data.Info.Files) > 5 {
			for _, v := range data.Info.Files {
				if len(v.Path) > 0 {
					t.Length += v.Length
					t.FileCount++
					tmpFiles = append(tmpFiles, file{
						Name:   v.Path[0],
						Length: v.Length,
					})
				}
			}
			sort.Sort(tmpFiles)
			if len(tmpFiles) >= 5 {
				t.Files = append(t.Files, tmpFiles[:5]...)
			} else {
				t.Files = append(t.Files, tmpFiles[:len(tmpFiles)]...)
			}
		} else {
			for _, v := range data.Info.Files {
				if len(v.Path) > 0 {
					t.Length += v.Length
					t.FileCount++
					t.Files = append(t.Files, file{
						Name:   v.Path[0],
						Length: v.Length,
					})
				}
			}
		}
	}

	b, _ := json.Marshal(t)
	if len(string(b)) > 1024 {
		return
	}
	err = insertData(t.Infohash, string(b))
	return
}

func insertData(hash string, content string) (err error) {
	switch hash[0] {
	case '0':
		var data torrent.Infohash0
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case '1':
		var data torrent.Infohash1
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case '2':
		var data torrent.Infohash2
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case '3':
		var data torrent.Infohash3
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case '4':
		var data torrent.Infohash4
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case '5':
		var data torrent.Infohash5
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case '6':
		var data torrent.Infohash6
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case '7':
		var data torrent.Infohash7
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case '8':
		var data torrent.Infohash8
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case '9':
		var data torrent.Infohash9
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case 'A':
		var data torrent.Infohasha
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case 'B':
		var data torrent.Infohashb
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case 'C':
		var data torrent.Infohashc
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case 'D':
		var data torrent.Infohashd
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case 'E':
		var data torrent.Infohashe
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	case 'F':
		var data torrent.Infohashf
		data.Infohash = hash
		data.Data = content
		_, err = utils.Config.Engine.Insert(&data)
	}
	return
}
