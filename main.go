package main

import (
	"fmt"
	"time"

	"github.com/gorhill/cronexpr"
)

func main() {
	nextTime := cronexpr.MustParse("* * * * *").Next(time.Now())
	fmt.Println("nextTime is ", nextTime)
}
