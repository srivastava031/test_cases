package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var wg sync.WaitGroup
var counter int64

func main() {
	for outerloop := 0; outerloop < 1000; outerloop++ {
		counter = 0
		wg.Add(20)
		for j := 0; j < 20; j++ {
			go func() {
				for i := 0; i < 20; i++ {
					atomic.AddInt64(&counter, 1)
					fmt.Println(atomic.LoadInt64(&counter))
				}
				wg.Done()
			}()

		}
		wg.Wait()
		fmt.Println("final value of counter", counter)

	}
}
