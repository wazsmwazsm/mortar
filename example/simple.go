package main

import (
	"context"
	"fmt"
	"github.com/wazsmwazsm/mortar"
	"sync"
)

func main() {
	// create a task pool with cap 10 (max 10 workers)
	pool, err := mortar.NewPool(10)
	if err != nil {
		panic(err)
	}

	wg := new(sync.WaitGroup)
	// create task
	task := &mortar.Task{
		Handler: func(v ...interface{}) {
			wg.Done()
			fmt.Println(v)
		},
	}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		// add parmas for task
		task.Params = []interface{}{i, i * 2, "hello"}
		// put task to pool
		pool.Put(context.Background(), task)
	}

	wg.Add(1)
	pool.Put(context.Background(), &mortar.Task{
		Handler: func(v ...interface{}) {
			wg.Done()
			fmt.Println(v)
		},
		Params: []interface{}{"hi!"}, // set params when create task
	})

	wg.Wait()

	// close pool graceful
	pool.Close()
	// put return error when pool already closed
	err = pool.Put(context.Background(), &mortar.Task{
		Handler: func(v ...interface{}) {},
	})
	if err != nil {
		fmt.Println(err) // pool already closed
	}
}
