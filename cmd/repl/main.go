package main

import (
	"bufio"
	"fmt"
	"os"
	// buffer_manager "github.com/roackb2/simple_db/internal/buffer"
)

func print_usage() {
	fmt.Println("Simple DB 0.0.1")
	fmt.Println("Type .exit to exit")
}

func print_prompt() {
	fmt.Print("db > ")
}

func main() {
	// buffer_manager.Allocate()
	print_usage()
	for {
		print_prompt()
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		if text == ".exit\n" {
			break
		}
		fmt.Println(text)
	}
}
