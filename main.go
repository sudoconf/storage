package main

import (
	"github.com/btlike/storage/crawl"
	"github.com/btlike/storage/utils"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	utils.Init()
	crawl.Run()
}
