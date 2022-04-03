package main

import (
	"fmt"

	"github.com/schollz/raw/src/sampswap"
)

func main() {
	err := sampswap.App()
	if err != nil {
		fmt.Println(err)
	}
}
