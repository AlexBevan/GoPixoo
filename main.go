package main

import (
	"os"

	"github.com/alexbevan/gopixoo/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
