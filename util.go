package main

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func isNumeric(str string) bool {
	for _, r := range str {
		if (r < '0' || r > '9') && r != '.' {
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
	for k := range operators {
		if k == str {
			return true
		}
	}
	return false
}

func indexOfSA(element string, array []string, startIndex int) int {
	for i := startIndex; i < len(array); i++ {
		if array[i] == element {
			return i
		}
	}
	return -1
}
