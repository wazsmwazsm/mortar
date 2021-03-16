package main

import (
	"context"
	"fmt"
	"time"

	"github.com/wazsmwazsm/mortar"
)

func main() {
	ctx := context.TODO()
	pool, err := mortar.NewPool(ctx, 10)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 20; i++ {
		pool.Put(&mortar.Task{
			Ctx: ctx,
			Handler: func(ctx context.Context, v ...interface{}) {
				fmt.Println(v)
			},
			Params: []interface{}{i},
		})
	}

	time.Sleep(1e9)
}
