package main

import (
	"context"
	"fmt"

	"github.com/wazsmwazsm/mortar"
)

func main() {
	case1()
}

func case1() {
	fmt.Println("--- case1 start ---")
	defer fmt.Println("--- case1 stoped ---")
	ctx := context.TODO()
	pool, err := mortar.NewPool(ctx, 1)
	if err != nil {
		panic(err)
	}

	pool.Put(&mortar.Task{
		Ctx: ctx,
		Handler: func(ctx context.Context, v ...interface{}) {
			panic("aaaaaa!")
		},
		Params: []interface{}{"hi!"},
	})

	pool.Put(&mortar.Task{
		Ctx: ctx,
		Handler: func(ctx context.Context, v ...interface{}) {
			fmt.Println(v)
		},
		Params: []interface{}{"hi!"},
	})

	pool.Close()
	err = pool.Put(&mortar.Task{
		Ctx:     ctx,
		Handler: func(ctx context.Context, v ...interface{}) {},
	})
	if err != nil {
		fmt.Println(err)
	}
}
