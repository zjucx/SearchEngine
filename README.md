## Search Engine by Golang
### Introduction
使用golang写的最简搜索引擎，目前主要分为4个模块:爬虫构造数据(可有可无，构建好自己的博客后删除此模块)、分词模块(用于对源文件进行分词，构建用于索引的tag，可改进优化)、倒排索引(对索引tag构建倒排索引)和BTree(实现搜索的数据结构，可与BPlusTree、RBTree等进行性能对比分析)。工程目录结构如下:
```
├── btree
│   ├── bplustree.go     ------btree
│   ├── internode.go     ------btree内节点的定义和操作
│   ├── leafnode.go      ------btree叶节点的定义和操作
│   └── node.go          ------节点的接口类
├── serment    
│   └── segment.go       ------分词部分，对爬虫爬取的数据进行分词
├── invertidx    
│   ├── dict.go          ------由分词部分构成的字典,实现词和整型数据的映射
│   ├── file.go          ------文件相关的操作
│   ├── index.go         ------倒排索引的主要实现,使用外排序算法对分词进行排序
│   └── page.go          ------暂未使用,文件映射功能
├── main.go
└── crawler
    ├── sinaLogin.go     ------模拟登陆模块，工程中实现了新浪微博的模拟登陆
    ├── crawler.go       ------爬虫模块的入口，将接口暴漏于分布式模块
    ├── request.go       ------包装请求
    └── scrawler.go      ------爬虫的辅助类
```
### Requirements
```
1. Golang(1.6)           -------开发语言
2. BTree、倒排索引、外排序等算法
```

### Implement
#### BTree
原理参考[BTree和B+Tree详解](http://blog.csdn.net/endlu/article/details/51720299)
算法实现分拆为内节点、叶结点和根结点三个文件。内节点和叶结点包含插入搜索和查找操作，具体实现可参考代码
```
写好了实现的大体思路，假期结束还未调通，有时间继续完善。
```
#### 倒排索引
使用外排序对各个文本的分词结果进行排序，生成一个倒排索引文件(目前只是单机实现)
```
输入为分词部分构造的词典，输出为设计好的倒排索引文件也是BTree的输入
```
#### 爬虫部分
构造原始数据
```
输出纯文本文件，是分词部分的输入
```
#### 分词
暂时使用结巴分词进行分词(改进使用深度学习进行分词)
```
输入为爬虫部分构造的数据，输出分词字典为倒排的输入
```
### Using
```
各个模块在main.go中独立运行、调试
```

### To Do List
```
1) BPlusTree、RBTree实现、性能分析
2) 应用个人博客的全文索引
3) 使用MapReduce进行外排序的分布式实现
4) 使用深度学习进行分词改进(测试发现结巴分词准确率略低，深度学习在NLP已经很成熟)
```
### Discussing
- [submit issue](https://github.com/zjucx/SearchEngineByGolang/issues/new)
- email: 862575451@qq.com
