package crawl

import (
	"net/http"
	"time"
)

type crawl struct {
	xunleiClient   *http.Client
	torcacheClient *http.Client
}

func newCrawl() (c crawl) {
	c.xunleiClient = &http.Client{
		Timeout: 10 * time.Second,
	}
	c.torcacheClient = &http.Client{
		Timeout: 5 * time.Second,
	}
	return
}
