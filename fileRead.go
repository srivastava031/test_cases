package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"unicode"
)

func main() {

	//getFile := flag.String("fileName", "file.txt", "file containing the parameters")

	data, _ := ioutil.ReadFile("file.txt")
	CdrCount := string(data)
	fmt.Println(CdrCount)
	CdrCountArray := strings.FieldsFunc(CdrCount, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})
	fmt.Println(CdrCountArray)

	fmt.Printf("\n%T\n", CdrCountArray[2])

	var nOofCdrGen, timeDelay []int
	for index, Cdrfile := range CdrCountArray {
		if index%2 == 0 {
			arrayVal, _ := strconv.Atoi(Cdrfile)
			nOofCdrGen = append(nOofCdrGen, arrayVal)

		} else {
			arrayVal, _ := strconv.Atoi(Cdrfile)
			timeDelay = append(timeDelay, arrayVal)
		}
	}
	fmt.Printf("===================================>nOofCdrGen= %T\n", nOofCdrGen)
	fmt.Println("nOofCdrGen=", nOofCdrGen)
	fmt.Println("timeDelay=", timeDelay)
	// for i := 0; i < 10; i++ {
	// 	fmt.Println("this printed", i)
	// }

	// fmt.Println()
	// fmt.Println()
	// fmt.Println()
	// //fmt.Println("timeDelay=", timeDelay, index)
	// //fmt.Println("nOofCdrGen", nOofCdrGen, index)

}

// fmt.Println()

// fmt.Println("[0]", CdrCountArray[0])
// fmt.Println("[1]", CdrCountArray[1])
// fmt.Println("[2]", CdrCountArray[2])
// fmt.Println("[3]", CdrCountArray[3])
// fmt.Println("[4]", CdrCountArray[4])
//}
