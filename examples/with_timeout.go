package main

import (
	"context"
	"fmt"
	"github.com/wazsmwazsm/mortar"
	"time"
)

func main() {
	pool, err := mortar.NewPool(10)
	if err != nil {
		panic(err)
	}

	task := &mortar.Task{
		Handler: func(v ...interface{}) {
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	for i := 0; i < 5; i++ {
		pool.Put(ctx, task) // 5 秒后所有 worker 自动销毁
	}

	fmt.Println(pool.GetRunningWorkers()) // 5

	time.Sleep(time.Second * 6)

	fmt.Println(pool.GetRunningWorkers()) // 0

}
