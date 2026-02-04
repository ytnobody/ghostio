package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 || os.Args[1] != "watch" {
		fmt.Fprintln(os.Stderr, "Usage: ghostio watch owner/repo")
		os.Exit(1)
	}
	fmt.Printf("Watching %s...\n", os.Args[2])
}
