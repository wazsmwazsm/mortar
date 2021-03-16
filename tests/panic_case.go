package main

import (
	"fmt"

	"github.com/wazsmwazsm/mortar"
)

func main() {
	case1()
}

func case1() {
	fmt.Println("--- case1 start ---")
	defer fmt.Println("--- case1 stoped ---")

	pool, err := mortar.NewPool(1)
	if err != nil {
		panic(err)
	}

	pool.Put(&mortar.Task{
		Handler: func(v ...interface{}) {
			panic("aaaaaa!")
		},
		Params: []interface{}{"hi!"},
	})

	pool.Put(&mortar.Task{
		Handler: func(v ...interface{}) {
			fmt.Println(v)
		},
		Params: []interface{}{"hi!"},
	})

	pool.Close()
	err = pool.Put(&mortar.Task{
		Handler: func(v ...interface{}) {},
	})
	if err != nil {
		fmt.Println(err)
	}
}
