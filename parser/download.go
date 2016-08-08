package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	// "io"
	// "io/ioutil"
	"net/http"

	// "github.com/btlike/storage/utils"
)

//define address
const (
	Xunlei   = "http://bt.box.n0808.com/%s/%s/%s.torrent"
	Torcache = "https://torcache.net/torrent/%s.torrent"
)

//define errors
var (
	ErrNotFound = errors.New("not found")
)

//DownloadXunlei torrent
func DownloadXunlei(hash string, client *http.Client) (mi MetaInfo, err error) {
	mi.InfoHash = hash
	if len(hash) != 40 {
		err = errors.New("invalid hash len")
		return
	}

	//从迅雷种子库查找
	address := fmt.Sprintf(Xunlei, hash[:2], hash[len(hash)-2:], hash)
	req0, err := http.NewRequest("GET", address, nil)
	if err != nil {
		return
	}
	req0.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req0)
	if err != nil {
		return
	}
	if resp != nil {
		defer func() {
			// io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		}()

		if resp.StatusCode == 200 {
			//解析种子
			err = mi.Parse(resp.Body)
		} else if resp.StatusCode == 404 {
			err = ErrNotFound
		} else {
			err = errors.New("refuse error")
		}
	}
	return
}

//DownloadTorcache torrent
func DownloadTorcache(hash string, client *http.Client) (mi MetaInfo, err error) {
	mi.InfoHash = hash
	if len(hash) != 40 {
		err = errors.New("invalid hash len")
		return
	}
	address := fmt.Sprintf(Torcache, hash)
	var req *http.Request
	req, err = http.NewRequest("GET", address, nil)
	if err != nil {
		return
	}
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36")

	// bg := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	// fmt.Println(time.Now().Sub(bg).String())

	defer func() {
		// io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		if resp.StatusCode == 404 {
			err = ErrNotFound
			return
		}
		err = errors.New("refuse error")
		return
	}

	//解析种子
	err = mi.Parse(resp.Body)
	if err != nil {
		err = fmt.Errorf("file cracked on torcache: %s,error is %s\n ", address, err)
		return
	}
	return
}

func pretty(v interface{}) {
	b, _ := json.MarshalIndent(v, " ", "  ")
	fmt.Println(string(b))
}
