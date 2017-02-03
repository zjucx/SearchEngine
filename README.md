## Search Engine by Golang

[![Build Status](https://travis-ci.org/zjucx/golang-webserver.svg?branch=master
)](http://120.27.39.169:8080/home)
[![Yii2](https://img.shields.io/badge/PoweredBy-ZjuCx-brightgreen.svg?style=flat)](http://120.27.39.169:8080/home)

### Introduction
使用golang开发的[分布式爬虫系统](https://github.com/zjucx/DistributedCrawler.git)，主要分为3个模块:[分布式框架](src/docs/framework.md)、[数据管理](src/docs/model.md)和[爬虫部分](src/docs/scrawler.md)。目录结构如下:
```
├── btree
│   ├── bplustree.go     ------btree
│   ├── internode.go     ------btree内节点的定义和操作
│   ├── leafnode.go      ------btree叶节点的定义和操作
│   └── node.go          ------节点的接口类
├── serment    
│   └── segment.go     ------分词部分，对爬虫爬取的数据进行分词
├── invertidx    
│   ├── dict.go       ------由分词部分构成的字典,实现词和整型数据的映射
│   ├── file.go       ------文件相关的操作
│   ├── index.go      ------倒排索引的主要实现,使用外排序算法对分词进行排序
│   └── page.go       ------暂未使用,文件映射功能
├── main.go
└── crawler
    ├── sinaLogin.go   ------模拟登陆模块，工程中实现了新浪微博的模拟登陆
    ├── crawler.go     ------爬虫模块的入口，将接口暴漏于分布式模块
    ├── request.go     ------包装请求
    └── scrawler.go    ------爬虫的辅助类
```
### Requirements
```
1. Docker(1.1x)   -------部署mongodb服务
2. Golang(1.6)    -------开发语言
3. Mongodb        -------持久化介质
4. Redis          -------优先级队列
```

### Implement
#### [分布式框架](src/docs/framework.md)
```
分为master节点和worker节点，master节点用于分发任务，worker节点用于任务执行。
```
#### [数据管理](src/docs/model.md)
```
分为持久化mongodb和内存数据库redis(实现优先级队列)。
```
#### [爬虫部分](src/docs/scrawler.md)
```
模拟登陆部分获取cookie，数据爬取部分。
```
### Using
```
<!--  Prepare redis servre and containers for worker  --!>
git clone https://github.com/zjucx/DistributedCrawler.git
cd DistributedCrawler
go get (代理代理代理)
// for master
go run main.go master masterip:port
// for workers
go run main.go worker masterip:port workerip:port
```

### To Do List
```
1) 爬虫系统的[web界面]()
2) 日志管理，可维可测功能
3) 使用Zookeeper实现分布式配置管理
3) 爬虫的单机操作
```
### Discussing
- [submit issue](https://github.com/zjucx/DistributedCrawler/issues/new)
- email: 862575451@qq.com
