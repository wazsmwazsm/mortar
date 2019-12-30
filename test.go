package main

import (
	"./mortar"
	"context"
	"fmt"
	"time"
)

func main() {
	pool, err := mortar.NewPool(10)
	if err != nil {
		panic(err)
	}

	// ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	// defer cancel()

	pool.Put(context.Background(), &mortar.Task{
		Handler: func(v ...interface{}) {
			fmt.Println(v)
		},
	})

	for i := 0; i < 10000; i++ {
		pool.Put(context.Background(), &mortar.Task{
			Handler: func(v ...interface{}) {
				time.Sleep(1e9)
				fmt.Println(v)
			},
			Params: []interface{}{i, 2, "hello"},
		})

		if i == 150 {
			pool.Close()
		}
	}

	for {
	}
}
