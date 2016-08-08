## Storage

根据infohash从xunlei和torcache获取torrent metadata存储到数据库，同时增加搜索引擎全文索引



## 特性

- 分表存储，设计容量6千万到8千万，均分在16张表中
- 丢弃torrent部分字段，节省90%网络流量
- 引擎健康检查，全部拒绝服务时，暂停抓取
- 多线程抓取
