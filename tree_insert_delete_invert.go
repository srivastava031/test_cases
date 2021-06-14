package main

import (
	"fmt"
)

type node struct {
	data  int
	left  *node
	right *node
}

func insert(root *node, val int) *node {
	if root == nil {
		root = getNewNode(val)
		return root
	} else if val <= root.data {
		root.left = insert(root.left, val)
	} else {
		root.right = insert(root.right, val)
	}
	return root
}

func getNewNode(val int) *node {
	newNode := new(node)
	newNode.data = val
	newNode.left, newNode.right = nil, nil
	return newNode
}

func display(n *node) {
	if n == nil {
		return
	}
	display(n.left)
	fmt.Printf("-->%d", n.data)
	display(n.right)
}
func del(n *node, val int) *node {
	if n == nil {
		return nil
	} else if val < n.data {
		n.left = del(n.left, val)
	} else if val > n.data {
		n.right = del(n.right, val)
	} else {
		if n.left == nil && n.right == nil {
			n = nil
			return n
		} else if n.left == nil {
			temp := n.right
			n.right = nil
			return temp
		} else if n.right == nil {
			temp := n.left
			n.left = nil
			return temp
		} else {
			findMin := func(n *node) *node {
				for n.left != nil {
					n = n.left
				}
				return n
			}
			minNode := findMin(n.right)
			n.data = minNode.data
			n.right = del(n.right, minNode.data)
		}
	}
	return n
}

func invert(root *node) {
	if root != nil {
		invert(root.left)
		invert(root.right)
		temp := root.left
		root.left = root.right
		root.right = temp
	}
	return
}

func main() {
	root := &node{}
	root = insert(root, 12)
	root = insert(root, 15)
	root = insert(root, 7)
	root = insert(root, 3)
	root = insert(root, 13)
	root = insert(root, 17)
	root = insert(root, 1)
	display(root)
	root = del(root, 17)
	fmt.Println()
	display(root)
	invert(root)
	fmt.Println()
	display(root)
}
