package logger

import "fmt"

const DEBUG = true

func Debug(msg string, args ...interface{}) {
	if DEBUG {
		fmt.Printf(msg, args...)
	}
}

func Error(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
}
