package test

import (
	"fmt"
	learn2 "helloDB/src/learn"
	"time"
)

func Mutex() {
	c := learn2.SafeCounter{V: make(map[string]int)}
	for i := 0; i < 1000; i++ {
		go c.Inc("hello")
	}

	time.Sleep(time.Second)
	fmt.Println(c.Value("hello"))
}
