package learn

import (
	"fmt"
	"golang.org/x/tour/tree"
)

func Walk(t *tree.Tree, ch chan int)  {
	sendVal(t, ch)
	close(ch)
}

func sendVal(t *tree.Tree, ch chan int)  {
	if t != nil {
		sendVal(t.Left, ch)
		ch <- t.Value
		sendVal(t.Right, ch)
	}
}

func IsSame(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go Walk(t1, ch1)
	go Walk(t2, ch2)
	for i := range ch1 {
		j := <- ch2
		fmt.Println("ch1: ", i, "	ch2: ", j)
		if i != j {
			return false
		}
	}
	return true
}



































