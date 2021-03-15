# Mortar
A goroutine task pool

一个简单好用的高性能任务池, 代码只有 100 多行。

## 版本更新日志

### v1.x
#### v1.1
- 部分冗余逻辑优化
#### v1.2
- 修复数据竞争 bug
#### v1.3
- 安全运行 worker
#### v1.4
- 退出等待 taskC 清空时增加 sleep 减少 cpu 负载
#### v1.5
- 优化锁，解决只有一个 woker 时产生 panic 后无发消费 task 导致 deadlock 的问题 (见 [issue 极端情况 #4](https://github.com/wazsmwazsm/mortar/issues/4))

## 解决什么问题

go 的 goroutine 提供了一种较线程而言更廉价的方式处理并发场景, 但 goroutine 太多会导致调度性能下降、GC 频繁、内存暴涨, 引发一系列问题。

mortar 限制了最多可启动的 goroutine 数量, 同时保持和原生 goroutine 一样的性能（在海量 goroutine 场景优于原生 goroutine）, 避免了上述问题。

## 原理

创建一个容量为 N 的池, 在池容量未满时, 每塞入一个任务（生产任务）, 任务池开启一个 worker (建立协程) 去处理任务（消费任务）。
当任务池容量赛满，每塞入一个任务（生产任务）, 任务会被已有的 N 个 worker 抢占执行（消费任务），达到协程限制的功能。

### 生产消费模型

队列: channel

生产任务: 将任务写入 channel

消费任务: worker（goroutine）从 channel 中读出任务执行

## 使用

### task struct

每个任务是一个结构体, Handler 是要执行的任务函数, Params 是要传入 Handler 的参数

```go
type Task struct {
	Handler func(v ...interface{})
	Params  []interface{}
}
```

### NewPool

NewPool() 方法创建一个任务池结构, 返回其指针

```go
func NewPool(capacity uint64) (*Pool, error)
```

### Put

Put() 方法来将一个任务放入池中, 如果任务池未满, 则启动一个 worker。

```go
func (p *Pool) Put(task *Task) error 
```

### GetCap

获取任务池容量, 创建任务池时已确定
```go
func (p *Pool) GetCap() uint64
```

### GetRunningWorkers

获取当前运行 worker 的数量
```go
func (p *Pool) GetRunningWorkers() uint64 
```

### Close()

安全关闭任务池。Close() 方法会先阻止 Put() 方法继续放入任务, 等待所有任务都被消费运行后, 销毁所有 worker, 关闭任务 channel。
```go
func (p *Pool) Close() 
```

### panic handler

每个 worker 都是一个原生 goroutine, 为保证程序的安全运行, 任务池会 recover 所有 worker 中的 panic, 并提供自定义的 panic 处理能力（不设置 PanicHandler 默认会打印 panic 的异常栈）。

```go
pool.PanicHandler = func(r interface{}) {
	// handle panic
	log.Println(r) 
}
```

## 例子

```go
package main

import (
	"fmt"
	"github.com/wazsmwazsm/mortar"
	"sync"
)

func main() {
	// 创建容量为 10 的任务池
	pool, err := mortar.NewPool(10)
	if err != nil {
		panic(err)
	}

	wg := new(sync.WaitGroup)

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		// 创建任务
		task := &mortar.Task{
			Handler: func(v ...interface{}) {
				wg.Done()
				fmt.Println(v)
			},
		}
		// 添加任务函数的参数
		task.Params = []interface{}{i, i * 2, "hello"}
		// 将任务放入任务池
		pool.Put(task)
	}

	wg.Add(1)
	// 再创建一个任务
	pool.Put(&mortar.Task{
		Handler: func(v ...interface{}) {
			wg.Done()
			fmt.Println(v)
		},
		Params: []interface{}{"hi!"}, // 也可以在创建任务时设置参数
	})

	wg.Wait()

	// 安全关闭任务池（保证已加入池中的任务被消费完）
	pool.Close()
	// 如果任务池已经关闭, Put() 方法会返回 ErrPoolAlreadyClosed 错误
	err = pool.Put(&mortar.Task{
		Handler: func(v ...interface{}) {},
	})
	if err != nil {
		fmt.Println(err) // print: pool already closed
	}
}

```

更多例子参考 examples 目录下的文件


## benchmark

100w 次执行，原子增量操作

模式 | 操作时间消耗 ns/op | 内存分配大小 B/op | 内存分配次数 allocs/op
-|-|-|-
原生 goroutine (100w goroutine) |	1596177880  |	103815552 	|  240022  
任务池开启 20 个 worker 20 goroutine) | 1378909099 	  | 15312 	  |    89 

### 对比

使用任务池和原生 goroutine 性能相近（略好于原生）

使用任务池比直接 goroutine 内存分配节省 7000 倍左右, 内存分配次数减少 2700 倍左右

> tips: 当任务为耗时任务时, 防止任务堆积（消费不过来）可以结合业务调整容量, 或根据业务控制每个任务的超时时间

# License

The Mortar is open-sourced software licensed under the [MIT license](http://opensource.org/licenses/MIT).
