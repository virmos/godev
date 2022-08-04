package main

import (
	"golang.org/x/tour/tree"
	"fmt"
)

func Walk(t *tree.Tree, channel chan int)  {
	var walk func(t *tree.Tree, channel chan int) 
	walk = func(t *tree.Tree, channel chan int) {
		if t != nil {
			walk(t.Left, channel)
			channel <- t.Value
			walk(t.Right, channel)	
		}
		
	}
	walk(t, channel)
	defer close(channel)

}

func Same(t1, t2 *tree.Tree) bool {
	ch1, ch2 := make(chan int), make(chan int)
	go Walk(t1, ch1)
	go Walk(t2, ch2)
	for {
		v1, ok1 := <-ch1
		v2, ok2 := <-ch2
		if !ok1 || !ok2 { // ch1 or ch2 close
			return ok1 == ok2
		}
		if v1 != v2 {
			return false
		}
	}
}

func test1() {
	fmt.Println(Same(tree.New(1), tree.New(1)))
}
