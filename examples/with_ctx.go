package main

import (
	"context"
	"fmt"
	"time"

	"github.com/wazsmwazsm/mortar"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := mortar.NewPool(10)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 20; i++ {
		pool.Put(&mortar.Task{
			Ctx: ctx, // use context
			Handler: func(ctx context.Context, v ...interface{}) {
				for {
					select {
					case <-ctx.Done():
						fmt.Printf("%v stoped\n", v) // after 5s print this
						return
					}
				}
			},
			Params: []interface{}{i},
		})
	}

	time.Sleep(10 * time.Second)
}
