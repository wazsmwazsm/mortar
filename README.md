# mortar
A goroutine task pool

一个简单好用的高性能任务池

## 解决什么问题

go 的 goroutine 提供了一种较线程而言更廉价的方式处理并发场景, 但 goroutine 太多会导致调度性能下降、GC 频繁、内存暴涨, 引发一系列问题。

mortar 限制了最多可启动的 goroutine 数量, 同时保持和原生 goroutine 一样的性能（在海量 goroutine 场景优于原生 goroutine）, 避免了上述问题。

## 原理

创建一个容量为 N 的池, 在池容量未满时, 每塞入一个任务（生产任务）, 任务池开启一个 worker (建立协程) 去处理任务（消费任务）。
当任务池容量赛满，每塞入一个任务（生产任务）, 任务会被已有的 N 个 worker 抢占执行，达到协程限制的功能。

### 生产消费模型

队列: channel

生产任务: 将任务写入 channel

消费任务: worker（goroutine）从 channel 中读出任务执行

