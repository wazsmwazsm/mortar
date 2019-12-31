package main

import (
	"context"
	"fmt"
	"github.com/wazsmwazsm/mortar"
	"sync"
)

func main() {
	pool, err := mortar.NewPool(10)
	if err != nil {
		panic(err)
	}

	task := &mortar.Task{
		Handler: func(v ...interface{}) {
			fmt.Println(v)
		},
		Params: []interface{}{i, i * 2, "hello"},
	}

	for i := 0; i < 1000; i++ {
		err := pool.Put(context.Background())
		if err != nil {
			fmt.Println(i, err)
		}

		if i == 550 {
			pool.Close()
		}
	}

	for {
	}
}
