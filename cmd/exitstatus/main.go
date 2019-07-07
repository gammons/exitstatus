package main

import (
	"fmt"
	"time"

	"github.com/gorhill/cronexpr"
)

func main() {
	nextTime := cronexpr.MustParse("0 0 29 2 *").Next(time.Now())
	fmt.Println("nextTime is ", nextTime)
}
