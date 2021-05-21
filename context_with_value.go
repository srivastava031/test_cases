package main

import (
	"fmt"
	"context"
	"time"
)

    type Values struct {
        m map[string]string
    }

func main() {
   vvv := Values{
        m: map[string]string{"1": "one"},
    }
	ctx := context.WithValue(context.Background(), "mv", vvv)
	go f(ctx)
	time.Sleep(3*time.Second)
}

func f(ctx context.Context){
	v := ctx.Value("mv")
        if v != nil {
                fmt.Println("found value:", v)
		}else{
			fmt.Println("value not found")
		}
} 
