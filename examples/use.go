package main

import (
	"fmt"
	"time"

	"github.com/wazsmwazsm/mortar"
)

func main() {
	pool, err := mortar.NewPool(10)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 20; i++ {
		pool.Put(&mortar.Task{
			Handler: func(v ...interface{}) {
				fmt.Println(v)
			},
			Params: []interface{}{i},
		})
	}

	time.Sleep(1e9)
}
