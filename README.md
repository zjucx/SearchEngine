## Search Engine by Golang

[![Build Status](https://travis-ci.org/zjucx/golang-webserver.svg?branch=master
)](http://120.27.39.169:8080/home)
[![Yii2](https://img.shields.io/badge/PoweredBy-ZjuCx-brightgreen.svg?style=flat)](http://120.27.39.169:8080/home)

### Introduction
使用golang开发的[分布式爬虫系统](https://github.com/zjucx/DistributedCrawler.git)，主要分为3个模块:[分布式框架](src/docs/framework.md)、[数据管理](src/docs/model.md)和[爬虫部分](src/docs/scrawler.md)。目录结构如下:
```
├── conf
│   └── app.conf       ------配置部分，数据库等信息的配置。还未开发。=。=
├── model    
│   ├── mongodb.go     ------爬虫的持久化介质，存储url和想要获取的数据
│   └── redismq.go     ------实用redis实现的优先级队列，master从mongodb获取url和向worker分发url
├── distribute    
│   ├── common.go      ------分布式系统的辅助类的定义等
│   ├── master.go      ------分布式系统的master节点，任务的分发调度
│   └── worker.go      ------分布式系统的worker节点，接受master的任务
├── main.go
└── scrawler           ------定义了数据库模型，用于与数据库交互
    ├── sinaLogin.go   ------模拟登陆模块，工程中实现了新浪微博的模拟登陆
    ├── scrawler.go    ------爬虫模块的入口，将接口暴漏于分布式模块
    ├── scheduler.go   ------爬虫的调度器，由于对master分发的url任务的预处理
    ├── downloader.go  ------爬虫的下载器，管理多个下载任务的同步等操作
    ├── spiders.go     ------爬虫的数据提取，用于提取resp的url和想要爬取的数据
    ├── pipeline.go    ------url和目的数据的持久化操作
    ├── request.go     ------封装的request请求
    └── utils.go       ------爬虫的辅助类
```
### Requirements
```
1. Docker(1.1x)   -------部署mongodb服务
2. Golang(1.6)    -------开发语言
3. Mongodb        -------持久化介质
4. Redis          -------优先级队列
```
### Screenshots
#### design
![](https://github.com/zjucx/redismq/blob/master/docs/distributeredis.bmp)
![](https://github.com/zjucx/redismq/blob/master/docs/contains.png)

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
