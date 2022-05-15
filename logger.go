package main

import (
	"fmt"
	"strconv"
)

func log(a []any, state string) {
	a = append([]any{"taylor: " + state + ":"}, a...)
	fmt.Println(a...)
}

func debug(a ...any) {
	if debugMode {
		log(a, "debug")
	}
}

func info(a ...any) {
	log(a, "info")
}

func warn(a ...any) {
	log(a, "warn")
}

func err(a ...any) {
	log(a, "error")
}

func errLine(message string, line int) {
	lineDenotion := " (" + strconv.Itoa(line) + ")"
	if line <= 0 {
		lineDenotion = ""
	}
	log([]any{message}, "error"+lineDenotion)
}
