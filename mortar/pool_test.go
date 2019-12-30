package mortar

import (
	"context"
	"sync/atomic"
	"testing"
)

var sum int64
var runTimes = 1000000

func demoTask(v ...interface{}) {
	for i := 0; i < 100; i++ {
		atomic.AddInt64(&sum, 1)
	}
}

func BenchmarkGoroutine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		go demoTask()
	}
}

func BenchmarkPut(b *testing.B) {
	pool, err := NewPool(10)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	task := &Task{
		Handler: demoTask,
	}

	for i := 0; i < b.N; i++ {
		pool.Put(ctx, task)
	}
}

func BenchmarkGoroutineSetTimes(b *testing.B) {

	for i := 0; i < runTimes; i++ {
		go demoTask()
	}
}

func BenchmarkPutSetTimes(b *testing.B) {
	pool, err := NewPool(10)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	task := &Task{
		Handler: demoTask,
	}

	for i := 0; i < runTimes; i++ {
		pool.Put(ctx, task)
	}
}
