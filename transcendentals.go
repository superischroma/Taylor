package main

import (
	"bufio"
	"math"
	"os"
	"strconv"
)

type Transcendental struct {
	arguments   []string
	description string
	operation   func([]string) (string, bool)
	variant     SymbolVariant
}

const (
	Radians = false
	Degrees = true
)

var trigMode = Radians

var transcendentals map[string]Transcendental

func radiansCheck(value float64) float64 {
	if trigMode {
		return (value * math.Pi) / 180
	}
	return value
}

func transSin(values []string) (string, bool) {
	return ftoa(math.Sin(radiansCheck(atof(values[0])))), true
}

func transCos(values []string) (string, bool) {
	return ftoa(math.Cos(radiansCheck(atof(values[0])))), true
}

func transTan(values []string) (string, bool) {
	return ftoa(math.Tan(radiansCheck(atof(values[0])))), true
}

func transCsc(values []string) (string, bool) {
	return ftoa(1.0 / math.Sin(radiansCheck(atof(values[0])))), true
}

func transSec(values []string) (string, bool) {
	return ftoa(1.0 / math.Cos(radiansCheck(atof(values[0])))), true
}

func transCot(values []string) (string, bool) {
	return ftoa(1.0 / math.Tan(radiansCheck(atof(values[0])))), true
}

func transSinh(values []string) (string, bool) {
	return ftoa(math.Sinh(radiansCheck(atof(values[0])))), true
}

func transCosh(values []string) (string, bool) {
	return ftoa(math.Cosh(radiansCheck(atof(values[0])))), true
}

func transTanh(values []string) (string, bool) {
	return ftoa(math.Tanh(radiansCheck(atof(values[0])))), true
}

func transCsch(values []string) (string, bool) {
	return ftoa(1 / math.Sinh(radiansCheck(atof(values[0])))), true
}

func transSech(values []string) (string, bool) {
	return ftoa(1 / math.Cosh(radiansCheck(atof(values[0])))), true
}

func transCoth(values []string) (string, bool) {
	return ftoa(1 / math.Tanh(radiansCheck(atof(values[0])))), true
}

func transArcsin(values []string) (string, bool) {
	return ftoa(math.Asin(radiansCheck(atof(values[0])))), true
}

func transArccos(values []string) (string, bool) {
	return ftoa(math.Acos(radiansCheck(atof(values[0])))), true
}

func transArctan(values []string) (string, bool) {
	return ftoa(math.Atan(radiansCheck(atof(values[0])))), true
}

func transArccsc(values []string) (string, bool) {
	return ftoa(math.Asin(1 / radiansCheck(atof(values[0])))), true
}

func transArcsec(values []string) (string, bool) {
	return ftoa(math.Acos(1 / radiansCheck(atof(values[0])))), true
}

func transArccot(values []string) (string, bool) {
	return ftoa(math.Atan(1 / radiansCheck(atof(values[0])))), true
}

func transArcsinh(values []string) (string, bool) {
	return ftoa(math.Asinh(radiansCheck(atof(values[0])))), true
}

func transArccosh(values []string) (string, bool) {
	return ftoa(math.Acosh(radiansCheck(atof(values[0])))), true
}

func transArctanh(values []string) (string, bool) {
	return ftoa(math.Atanh(radiansCheck(atof(values[0])))), true
}

func transArccsch(values []string) (string, bool) {
	return ftoa(math.Asinh(1 / radiansCheck(atof(values[0])))), true
}

func transArcsech(values []string) (string, bool) {
	return ftoa(math.Acosh(1 / radiansCheck(atof(values[0])))), true
}

func transArccoth(values []string) (string, bool) {
	return ftoa(math.Atanh(1 / radiansCheck(atof(values[0])))), true
}

func transLog(values []string) (string, bool) {
	return ftoa(math.Log10(atof(values[1])) / math.Log10(atof(values[0]))), true
}

func transLg(values []string) (string, bool) {
	return ftoa(math.Log10(atof(values[0]))), true
}

func transLn(values []string) (string, bool) {
	return ftoa(math.Log(atof(values[0]))), true
}

func transExp(values []string) (string, bool) {
	return ftoa(math.Exp(atof(values[0]))), true
}

func transSqrt(values []string) (string, bool) {
	return ftoa(math.Sqrt(atof(values[0]))), true
}

func transPi(values []string) (string, bool) {
	return ftoa(math.Pi), true
}

func transE(values []string) (string, bool) {
	return ftoa(math.E), true
}

func transRad(values []string) (string, bool) {
	return ftoa((atof(values[0]) * math.Pi) / 180.0), true
}

func transDeg(values []string) (string, bool) {
	return ftoa((180.0 * atof(values[0])) / math.Pi), true
}

func transAbs(values []string) (string, bool) {
	return ftoa(math.Abs(atof(values[0]))), true
}

func transInc(values []string) (string, bool) {
	filename := values[0]
	ok, _ := interpret(filename)
	if !ok {
		return "", false
	}
	return "", true
}

func transHirt(values []string) (string, bool) {
	return ftoa(math.Pow(atof(values[1]), 1.0/atof(values[0]))), true
}

func transRead(values []string) (string, bool) {
	reader := bufio.NewReader(os.Stdin)
	value, _ := reader.ReadString('\n')
	value = value[:len(value)-2] // remove irrelevant characters
	ok, fvalue := interpretLine(value, -1, false)
	if !ok {
		return "", false
	}
	_, e := strconv.ParseFloat(fvalue, 64)
	if e != nil {
		err("user-generated input is not a number")
		return "", false
	}
	return fvalue, true
}

func transExit(values []string) (string, bool) {
	os.Exit(0)
	return "", true
}
