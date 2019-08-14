package main

import (
	"fmt"
	"os"

	"go.fraixed.es/storj/sntools/dbcleaner/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
