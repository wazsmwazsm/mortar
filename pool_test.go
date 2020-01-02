package mortar

import (
	"sync"
	"sync/atomic"
	"testing"
)

var sum int64
var runTimes = 1000000

var wg = sync.WaitGroup{}

func demoTask(v ...interface{}) {
	for i := 0; i < 100; i++ {
		atomic.AddInt64(&sum, 1)
	}
}

func demoTask2(v ...interface{}) {
	defer wg.Done()
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
		b.Error(err)
	}

	task := &Task{
		Handler: demoTask,
	}

	for i := 0; i < b.N; i++ {
		pool.Put(task)
	}
}

func BenchmarkGoroutineTimelife(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go demoTask2()
	}
	wg.Wait()
}

func BenchmarkPutTimelife(b *testing.B) {
	pool, err := NewPool(10)
	if err != nil {
		b.Error(err)
	}

	task := &Task{
		Handler: demoTask2,
	}

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		pool.Put(task)
	}
	wg.Wait()

}

func BenchmarkGoroutineSetTimes(b *testing.B) {

	for i := 0; i < runTimes; i++ {
		go demoTask()
	}
}

func BenchmarkPoolPutSetTimes(b *testing.B) {
	pool, err := NewPool(20)
	if err != nil {
		b.Error(err)
	}

	task := &Task{
		Handler: demoTask,
	}

	for i := 0; i < runTimes; i++ {
		pool.Put(task)
	}
}

func BenchmarkGoroutineTimeLifeSetTimes(b *testing.B) {

	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		go demoTask2()
	}
	wg.Wait()
}

func BenchmarkPoolTimeLifeSetTimes(b *testing.B) {
	pool, err := NewPool(20)
	if err != nil {
		b.Error(err)
	}

	task := &Task{
		Handler: demoTask2,
	}

	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		pool.Put(task)
	}

	wg.Wait()
}
