package parser

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/btlike/database/torrent"
	"github.com/btlike/storage/utils"
)

func TestTorcache(t *testing.T) {
	engine, err := xorm.NewEngine("mysql", "root:yanyuanld2046@tcp(10.99.0.12:3306)/torrent?charset=utf8")
	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
		// Transport: &http.Transport{},
	}

	var (
		NotFound int64
		Refuse   int64
		Found    int64
	)

	go func() {
		for {
			fmt.Println("not found", NotFound, "refuse", Refuse, "found", Found)
			time.Sleep(time.Second * 10)
		}
	}()

	for {
		var pres []torrent.PreInfohash
		err := engine.Limit(100, 0).Find(&pres)
		if err != nil {
			utils.Log().Println(err)
			time.Sleep(time.Second * 10)
			continue
		}

		for _, v := range pres {
			go func(hash string) {
				err = parseTorcache(hash, client)
				if err != nil {
					if err == ErrNotFound {
						NotFound++
					} else {
						Refuse++
						fmt.Println(err)
					}
				} else {
					Found++
				}
			}(v.Infohash)
		}
		time.Sleep(time.Second)
		// utils.Log().Println(address2)
		// mi.DumpTorrentMetaInfo()
	}
}

func parseTorcache(hash string, client *http.Client) (err error) {
	address2 := fmt.Sprintf(Torcache, hash)
	var req *http.Request
	req, err = http.NewRequest("GET", address2, nil)
	if err != nil {
		return
	}
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36")

	// bg := time.Now()
	resp2, err := client.Do(req)
	if err != nil {
		return
	}
	// fmt.Println(time.Now().Sub(bg).String())

	defer func() {
		// io.Copy(ioutil.Discard, resp2.Body)
		resp2.Body.Close()
	}()

	if resp2.StatusCode != 200 {
		err = ErrNotFound
		return
	}

	//解析种子
	var mi MetaInfo
	err = mi.Parse(resp2.Body)
	if err != nil {
		err = fmt.Errorf("file cracked on torcache: %s,error is %s\n ", address2, err)
		return
	}
	return
}
