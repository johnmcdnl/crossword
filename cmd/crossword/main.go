package main

import (
	"github.com/johnmcdnl/crossword"
	"time"
	"fmt"
)

func main() {

	for i := 0; i <= 20; i++ {
		c := crossword.New(15)
		time.Sleep(500 * time.Millisecond)
		fmt.Println(c)
	}

}
