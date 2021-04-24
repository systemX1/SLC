package test

import (
	"fmt"
	"golang.org/x/tour/tree"
	learn2 "SLC/src/learn"
)

func Tree()  {
	fmt.Println("t1 and t2: ", learn2.IsSame(tree.New(1), tree.New(1)))
}
