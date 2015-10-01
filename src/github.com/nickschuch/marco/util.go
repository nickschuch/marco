package main

import (
	"fmt"
	"os"
)

func Check(e error) {
	if e != nil {
		Exit(e)
	}
}

func Exit(e error) {
	fmt.Println(e)
	os.Exit(1)
}
