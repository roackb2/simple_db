package main

import (
	"fmt"

	buffer_manager "github.com/roackb2/simple_db/internal/buffer"
)

func main() {
	fmt.Println("Running REPL")
	buffer_manager.Allocate()
}
