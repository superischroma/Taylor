package main

import (
	"strconv"
	"strings"
)

func isNumeric(str string) bool {
	if str == "-" {
		return false
	}
	for _, r := range str {
		if (r < '0' || r > '9') && r != '.' && r != '-' {
			return false
		}
	}
	return true
}

func isAlpha(str string) bool {
	for _, r := range str {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return false
		}
	}
	return true
}

func isOperator(str string) bool {
	_, exists := operators[str]
	return exists
}

func isStringLiteral(str string) bool {
	return len(str) >= 2 && (str[0] == '"' && str[len(str)-1] == '"') || (str[0] == '\'' && str[len(str)-1] == '\'')
}

func unwrapStringLiteral(str string) string {
	if !isStringLiteral(str) {
		return str
	}
	if len(str) == 2 {
		return ""
	}
	return str[1 : len(str)-1]
}

func indexOfSA(element string, array []string, startIndex int) int {
	for i := startIndex; i < len(array); i++ {
		if array[i] == element {
			return i
		}
	}
	return -1
}

func ftoa(value float64) string {
	return strconv.FormatFloat(value, 'E', -1, 64)
}

func atof(str string) float64 {
	value, _ := strconv.ParseFloat(str, 64)
	return value
}

func btois(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func itob(f float64) bool {
	return f != 0
}

func stringifyFunction(fSymbol *Symbol) string {
	str := "[U] " + fSymbol.name + "("
	for i, parameter := range fSymbol.children {
		if i != 0 {
			str += ", "
		}
		str += parameter.name
	}
	str += ")"
	return str
}

func stringifyTransFunction(name string, transcendental *Transcendental) string {
	return "[T] " + name + "(" + strings.Join(transcendental.arguments, ", ") + "): " + transcendental.description
}
