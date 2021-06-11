package main

import (
	"fmt"
)

type node struct {
	data int
	next *node
	prev *node
}
type list struct {
	head *node
	tail *node
}

func (l *list) insert(val int) {
	temp := &node{
		data: val,
		next: l.head,
	}
	if l.head != nil {
		l.head.prev = temp
	}
	l.head = temp

	ls := l.head
	for ls.next != nil {
		ls = ls.next
	}
	l.tail = ls
}
func (l *list) display() {
	ll := l.head
	for ll != nil {
		fmt.Printf("--> %d", ll.data)
		ll=ll.next
	}
}
func (l *list) del(n int) {
	ll := l.head
	for i := 0; i < n; i++ {
		ll = ll.next
	}
	ll.next.prev = ll.prev
	ll.prev.next = ll.next
}
func (l *list) deln(n int) {
	ll := l.head
	for ll != nil {
		if ll.data == n {
			break
		} else {
			ll = ll.next
		}
	}
	ll.next.prev = ll.prev
	ll.prev.next = ll.next
}
func (l *list) reverse() {
	currentnode := l.head
	for currentnode != nil {
	nextnode:=currentnode.next
	currentnode.next=currentnode.prev
	currentnode.prev=nextnode
	currentnode=nextnode
	}
	currentnode=l.head
	l.head=l.tail
	l.tail=currentnode

}

func main() {
	lis := &list{}
	lis.insert(4)
	lis.insert(5)
	lis.insert(6)
	lis.insert(7)
	lis.insert(8)
	lis.insert(9)
	lis.insert(10)
	lis.insert(11)
	lis.display()
	lis.del(3)
	fmt.Println()
	lis.display()
	lis.reverse()
	fmt.Println()
	lis.display()
	lis.deln(9)
	fmt.Println()
	lis.display()
}
