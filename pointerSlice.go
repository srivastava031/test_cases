package main

import "fmt"

func main() {
	var p []int //same for array var p =[]int{1,2,3}
	modifySlice(&p)
	fmt.Println(p)
}
func modifySlice(p *[]int) {
	k := *p
	fmt.Println("k := *p ==>", k)
	*p = append(k, 4, 5)
	fmt.Println(*p)
}

// result is =
// k := *p ==> []
// [4 5]
// [4 5]
/*
func main() {
	var s = []string{"1", "2", "3"}
	modifySlice(s)
	fmt.Println(s[0])
	fmt.Println(s)
}

func modifySlice(i []string) {
	i[0] = "3"
	i = append(i, "4")
}*/
//here the result will be like this
// Result:
// 3
// [3 2 3]
/*reason

So there are two things here, in the function modifySlice there’s two different access to the slice.
The slice is, by definition, a pointer to an underlying array.

This code i[0] = "3" takes the position 0 in i and set it to "3",
 which is the original slice, since even when it doesn’t seem like it’s a pointer, it still is.

Here’s the issue: the append function makes a check… Quoting the docs:

If it has sufficient capacity, the destination is resliced to accommodate the new elements.
If it does not, a new underlying array will be allocated.

So when you do i = append(i, "4") you’re essentially saying: "add ‘4’ to this slice,
but since the original “i” slice has reached its maximum capacity,
create a new one, add “4” and then set it to “i”. Since “i” only exists within the modifySlice() function,
you created a new slice that you’ll
never return or use anymore.

Each slice has a length and a capacity. “The length of a slice is the number of elements it contains.
The capacity of a slice is the number of elements in the underlying array,
counting from the first element in the slice.”  So when you declared var s = []string{"1", "2", "3"},
you created a slice with length of 3 and capacity of 3.
Changing one value that is already in the slice won’t create another slice, but appending to a slice that
doesn’t have the capacity to hold more items will create a new one.

*/
