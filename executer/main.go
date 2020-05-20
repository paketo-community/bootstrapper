package main

import (
	"fmt"

	"github.com/joshzarrabi/cnb-bootstrapper/bootstrapper"
)

func main() {
	err := bootstrapper.Bootstrap()
	if err != nil {
		fmt.Println(err)
	}
}
