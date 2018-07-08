package main

import (
	"github.com/sirupsen/logrus"
	"github.com/johnmcdnl/crossword"
	"time"
	"fmt"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)


	c := crossword.New(15)

	time.Sleep(500 * time.Millisecond)
	fmt.Println(c)
}
