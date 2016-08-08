package crawl

import (
	"time"

	"github.com/btlike/database/torrent"
	"github.com/btlike/storage/parser"
	"github.com/btlike/storage/utils"
)

//Run the spider
func Run() {
	Manager.run()
	Crawl()
}

//Crawl from xunlei ...
func Crawl() {
	worker := func(jobs <-chan string, resultChan chan<- string) {
		crawl := newCrawl()
		for infohash := range jobs {
			if !Manager.crawStatus[Xunlei].pauseCrawl ||
				!Manager.crawStatus[Torcache].pauseCrawl {
				//至少有一个引擎在服务时，直接删除即可，防止引擎都不服务时，疯狂删数据
				resultChan <- infohash
			}

			var data parser.MetaInfo
			var err error

			if !Manager.crawStatus[Xunlei].pauseCrawl && !Manager.crawStatus[Xunlei].stopCrawl {
				data, err = parser.DownloadXunlei(infohash, crawl.xunleiClient)
				if err != nil {
					if err == parser.ErrNotFound {
						Manager.crawStatus[Xunlei].notFoundCount++
					} else {
						Manager.crawStatus[Xunlei].refuseCount++
					}
					//报错，如果torcache在服务，进入torcache流程，否则进入下一循环
					if !Manager.crawStatus[Torcache].pauseCrawl {
						goto tocache
					} else {
						continue
					}
				} else {
					//没报错，进入存储流程
					goto store
				}
			}

		tocache:
			if !Manager.crawStatus[Torcache].pauseCrawl && !Manager.crawStatus[Torcache].stopCrawl {
				// bg := time.Now()
				data, err = parser.DownloadTorcache(infohash, crawl.torcacheClient)
				if err != nil {
					if err == parser.ErrNotFound {
						Manager.crawStatus[Torcache].notFoundCount++
					} else {
						Manager.crawStatus[Torcache].refuseCount++
					}
					//报错，进入下一循环
					continue
				}
			} else {
				continue
			}

			if Manager.crawStatus[Torcache].pauseCrawl &&
				Manager.crawStatus[Xunlei].pauseCrawl {
				//预防引擎都没有服务时，直接进入下一循环
				continue
			}

		store:
			err = Store(data)
			if err != nil {
				resultChan <- infohash
				continue
			}

			//全文索引
			err = IndexManager.add(data)
			if err != nil {
				utils.Log().Println(err)
			}

			resultChan <- infohash
			Manager.storeCount++
		}
	}

	jobChan := make(chan string, DownloadChanLength)
	resultChan := make(chan string, DownloadChanLength)
	defer close(resultChan)
	for i := 0; i < 300; i++ {
		go worker(jobChan, resultChan)
	}

	go func() {
		var infohashs []string
		for infohash := range resultChan {
			if len(infohash) == 40 {
				infohashs = append(infohashs, infohash)
				if len(infohashs) >= 100 {
					_, err := utils.Config.Engine.In("infohash", infohashs).Delete(&torrent.PreInfohash{})
					infohashs = make([]string, 0)
					if err != nil {
						utils.Log().Println("delete error", err)
					}
				}
			}
		}
	}()

	for {
		if (Manager.crawStatus[Xunlei].pauseCrawl || Manager.crawStatus[Xunlei].stopCrawl) &&
			Manager.crawStatus[Torcache].pauseCrawl {
			utils.Log().Println("全部引擎拒绝服务,暂停抓取,等待10分钟")
			time.Sleep(time.Minute * 10)
			Manager.crawStatus[Xunlei] = &crawStatus{}
			Manager.crawStatus[Torcache] = &crawStatus{}
		}

		var pres []torrent.PreInfohash
		err := utils.Config.Engine.OrderBy("id").Limit(1000, 0).Find(&pres)
		if err != nil {
			utils.Log().Println(err)
			time.Sleep(time.Second * 10)
			continue
		}

		if len(pres) == 0 {
			time.Sleep(time.Second * 60)
		}
		for _, v := range pres {
			jobChan <- v.Infohash
		}
	}
}
