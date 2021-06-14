package main

import (
	"fmt"
)

func maxHeap(a []int, n, i int) {
	largest := i
	l := (2 * i) + 1
	r := (2 * i) + 2
	if l < n && a[l] > a[largest] {
		largest = l
	}
	if r < n && a[r] > a[largest] {
		largest = r
	}
	if largest != i {
		temp := a[largest]
		a[largest] = a[i]
		a[i] = temp
		maxHeap(a, n, largest)
	}
}

func main() {
	p := []int{2,20,12,10,3,4,1,7,8,52 }
	n := len(p)
	for i := (n / 2) - 1; i >= 0; i-- {
		maxHeap(p, n, i)
	}
	fmt.Println(p)
	for i := n - 1; i >= 0; i-- {
		temp := p[0]
		p[0] = p[i]
		p[i] = temp
		fmt.Println("---->>",p)
		maxHeap(p, i, 0)
		fmt.Println(">>>>>>",p)
	}
	fmt.Println(p)

}
