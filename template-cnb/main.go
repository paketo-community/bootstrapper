package main

import "github.com/paketo-buildpacks/packit"

func main() {
	packit.Run(Detect(), Build())
}
