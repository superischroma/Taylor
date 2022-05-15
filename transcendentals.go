package main

import (
	"math"
)

type Transcendental struct {
	argCount  int
	operation func(...string) string
	variant   SymbolVariant
}

const (
	Radians = false
	Degrees = true
)

var trigMode = Radians

var transcendentals map[string]Transcendental

func radiansCheck(value float64) float64 {
	if trigMode == Degrees {
		return (value * math.Pi) / 180
	}
	return value
}

func transSin(values ...string) string {
	return ftoa(math.Sin(radiansCheck(atof(values[0]))))
}

func transCos(values ...string) string {
	return ftoa(math.Cos(radiansCheck(atof(values[0]))))
}

func transTan(values ...string) string {
	return ftoa(math.Tan(radiansCheck(atof(values[0]))))
}

func transCsc(values ...string) string {
	return ftoa(1.0 / math.Sin(radiansCheck(atof(values[0]))))
}

func transSec(values ...string) string {
	return ftoa(1.0 / math.Cos(radiansCheck(atof(values[0]))))
}

func transCot(values ...string) string {
	return ftoa(1.0 / math.Tan(radiansCheck(atof(values[0]))))
}

func transSinh(values ...string) string {
	return ftoa(math.Sinh(radiansCheck(atof(values[0]))))
}

func transCosh(values ...string) string {
	return ftoa(math.Cosh(radiansCheck(atof(values[0]))))
}

func transTanh(values ...string) string {
	return ftoa(math.Tanh(radiansCheck(atof(values[0]))))
}

func transCsch(values ...string) string {
	return ftoa(1 / math.Sinh(radiansCheck(atof(values[0]))))
}

func transSech(values ...string) string {
	return ftoa(1 / math.Cosh(radiansCheck(atof(values[0]))))
}

func transCoth(values ...string) string {
	return ftoa(1 / math.Tanh(radiansCheck(atof(values[0]))))
}

func transArcsin(values ...string) string {
	return ftoa(math.Asin(radiansCheck(atof(values[0]))))
}

func transArccos(values ...string) string {
	return ftoa(math.Acos(radiansCheck(atof(values[0]))))
}

func transArctan(values ...string) string {
	return ftoa(math.Atan(radiansCheck(atof(values[0]))))
}

func transArccsc(values ...string) string {
	return ftoa(1 / math.Asin(radiansCheck(atof(values[0]))))
}

func transArcsec(values ...string) string {
	return ftoa(1 / math.Acos(radiansCheck(atof(values[0]))))
}

func transArccot(values ...string) string {
	return ftoa(1 / math.Atan(radiansCheck(atof(values[0]))))
}

func transArcsinh(values ...string) string {
	return ftoa(math.Asinh(radiansCheck(atof(values[0]))))
}

func transArccosh(values ...string) string {
	return ftoa(math.Acosh(radiansCheck(atof(values[0]))))
}

func transArctanh(values ...string) string {
	return ftoa(math.Atanh(radiansCheck(atof(values[0]))))
}

func transArccsch(values ...string) string {
	return ftoa(1 / math.Asinh(radiansCheck(atof(values[0]))))
}

func transArcsech(values ...string) string {
	return ftoa(1 / math.Acosh(radiansCheck(atof(values[0]))))
}

func transArccoth(values ...string) string {
	return ftoa(1 / math.Atanh(radiansCheck(atof(values[0]))))
}

func transLog(values ...string) string {
	return ftoa(math.Log10(atof(values[1])) / math.Log10(atof(values[0])))
}

func transLg(values ...string) string {
	return ftoa(math.Log10(atof(values[0])))
}

func transLn(values ...string) string {
	return ftoa(math.Log(atof(values[0])))
}

func transExp(values ...string) string {
	return ftoa(math.Exp(atof(values[0])))
}

func transSqrt(values ...string) string {
	return ftoa(math.Sqrt(atof(values[0])))
}

func transPi(values ...string) string {
	return ftoa(math.Pi)
}

func transE(values ...string) string {
	return ftoa(math.E)
}

func transRad(values ...string) string {
	return ftoa((atof(values[0]) * math.Pi) / 180.0)
}

func transDeg(values ...string) string {
	return ftoa((180.0 * atof(values[0])) / math.Pi)
}

func transAbs(values ...string) string {
	return ftoa(math.Abs(atof(values[0])))
}

func transInc(values ...string) string {
	filename := values[0]
	ok, _ := interpret(filename)
	if !ok {
		panic(-1)
	}
	return ""
}

func transHirt(values ...string) string {
	return ftoa(math.Pow(atof(values[1]), 1.0/atof(values[0])))
}
