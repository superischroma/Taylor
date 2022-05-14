package main

import (
	"fmt"
	"strconv"
)

func log(message string, state string) {
	fmt.Printf("taylor: %s: %s\n", state, message)
}

func debug(message string) {
	log(message, "debug")
}

func info(message string) {
	log(message, "info")
}

func warn(message string) {
	log(message, "warn")
}

func err(message string) {
	log(message, "error")
}

func lineError(message string, line int) {
	log(message, "error ("+strconv.Itoa(line)+")")
}
