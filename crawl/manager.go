package crawl

import (
	"sync"
	"time"

	"github.com/btlike/storage/utils"
)

//const
const (
	DownloadChanLength = 256
	Xunlei             = "xunlei"
	Torcache           = "torcache"
)

//Manager spider
var Manager manager

func (p *manager) run() {
	p.initChan()
	p.crawStatus = make(map[string]*crawStatus)
	p.crawStatus[Xunlei] = &crawStatus{}
	p.crawStatus[Torcache] = &crawStatus{}
	go p.monitor()
}

func (p *manager) initChan() {}

type crawStatus struct {
	preNotFoundCount int64
	preRefuseCount   int64

	notFoundCount int64
	refuseCount   int64

	pauseCrawl bool
	stopCrawl  bool
	pauseTime  time.Time
}

type manager struct {
	storeCount int64
	crawStatus map[string]*crawStatus
	sync.Mutex
}

func (p *manager) monitor() {
	go func() {
		for {
			for k, v := range p.crawStatus {
				if k == Xunlei {
					if (v.notFoundCount-v.preNotFoundCount)/5 < (v.refuseCount - v.preRefuseCount) {
						v.pauseCrawl = true
						v.pauseTime = time.Now()
					}
				} else if k == Torcache {
					if (v.notFoundCount - v.preNotFoundCount) < (v.refuseCount - v.preRefuseCount) {
						v.pauseCrawl = true
						v.pauseTime = time.Now()
					}
				}
				v.preNotFoundCount = v.notFoundCount
				v.preRefuseCount = v.refuseCount

				//释放被暂停的服务
				if time.Now().Sub(v.pauseTime).Minutes() >= 10 {
					v.pauseCrawl = false
				}

				utils.Log().Printf("%s未找到数量:%v,拒绝数量:%v\n", k, v.notFoundCount, v.refuseCount)
			}
			utils.Log().Println("此次运行已存储数量:", p.storeCount)
			time.Sleep(time.Minute)
		}
	}()
}
