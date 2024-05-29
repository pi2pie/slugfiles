/*
Copyright © 2024 pi2pie,Ltd.
Licensed under the MIT License.
*/
package main

import (
	"fmt"
	"os"

	"github.com/pi2pie/slugfiles/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}