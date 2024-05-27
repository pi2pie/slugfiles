package main

import (
	"fmt"
	"os"
	"pi2pie/slugfiles-rename/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}