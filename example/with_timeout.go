package main

import (
	"context"
	"fmt"
	"github.com/wazsmwazsm/mortar"
	"time"
)

func main() {
	// create a task pool with cap 10 (max 10 workers)
	pool, err := mortar.NewPool(10)
	if err != nil {
		panic(err)
	}

	// create task
	task := &mortar.Task{
		Handler: func(v ...interface{}) {
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	for i := 0; i < 5; i++ {
		// put task to pool
		pool.Put(ctx, task) // 5 second later workers will be destroyed
	}

	fmt.Println(pool.GetRunningWorkers()) // 5

	time.Sleep(time.Second * 6)

	fmt.Println(pool.GetRunningWorkers()) // 0

}
