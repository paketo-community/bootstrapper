package main

import (
	"fmt"

	"github.com/joshzarrabi/template-cnb/bootstrapper"
)

func main() {
	err := bootstrapper.Bootstrap()
	if err != nil {
		fmt.Println(err)
	}
}
