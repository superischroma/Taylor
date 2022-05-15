package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type LineResult int

const (
	outputLineResult   LineResult = 0
	funcDefLineResult  LineResult = 1
	constDefLineResult LineResult = 2
)

var debugMode = false

var operators = map[string]int{
	"=":  0,
	"<":  0,
	">":  0,
	"<=": 0,
	">=": 0,
	"+":  1,
	"-":  1,
	"*":  2,
	"/":  2,
	"^":  3,
}

var tokenizerRegExp *regexp.Regexp
var symbols map[string]*Symbol

func main() {
	files := []string{}
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "-debug" {
			debugMode = true
		} else {
			files = append(files, arg)
		}
	}
	if len(files) == 0 {
		err("no input files")
		return
	}
	transcendentals = map[string]Transcendental{
		"sin":     Transcendental{1, transSin, functionVariant},
		"cos":     Transcendental{1, transCos, functionVariant},
		"tan":     Transcendental{1, transTan, functionVariant},
		"csc":     Transcendental{1, transCsc, functionVariant},
		"sec":     Transcendental{1, transSec, functionVariant},
		"cot":     Transcendental{1, transCot, functionVariant},
		"sinh":    Transcendental{1, transSinh, functionVariant},
		"cosh":    Transcendental{1, transCosh, functionVariant},
		"tanh":    Transcendental{1, transTanh, functionVariant},
		"csch":    Transcendental{1, transCsch, functionVariant},
		"sech":    Transcendental{1, transSech, functionVariant},
		"coth":    Transcendental{1, transCoth, functionVariant},
		"arcsin":  Transcendental{1, transArcsin, functionVariant},
		"arccos":  Transcendental{1, transArccos, functionVariant},
		"arctan":  Transcendental{1, transArctan, functionVariant},
		"arccsc":  Transcendental{1, transArccsc, functionVariant},
		"arcsec":  Transcendental{1, transArcsec, functionVariant},
		"arccot":  Transcendental{1, transArccot, functionVariant},
		"arcsinh": Transcendental{1, transArcsinh, functionVariant},
		"arccosh": Transcendental{1, transArccosh, functionVariant},
		"arctanh": Transcendental{1, transArctanh, functionVariant},
		"arccsch": Transcendental{1, transArccsch, functionVariant},
		"arcsech": Transcendental{1, transArcsech, functionVariant},
		"arccoth": Transcendental{1, transArccoth, functionVariant},
		"log":     Transcendental{2, transLog, functionVariant},
		"lg":      Transcendental{1, transLg, functionVariant},
		"ln":      Transcendental{1, transLn, functionVariant},
		"exp":     Transcendental{1, transExp, functionVariant},
		// determinant when?
		"sqrt": Transcendental{1, transSqrt, functionVariant},
		"pi":   Transcendental{0, transPi, constantVariant},
		"e":    Transcendental{0, transE, constantVariant},
		"rad":  Transcendental{1, transRad, functionVariant},
		"deg":  Transcendental{1, transDeg, functionVariant},
		"abs":  Transcendental{1, transAbs, functionVariant},
		"inc":  Transcendental{1, transInc, functionVariant},
		"hirt": Transcendental{2, transHirt, functionVariant},
	}
	tokenizerRegExp, _ = regexp.Compile("radians|degrees|sinh|cosh|tanh|csch|sech|coth|sin|cos|tan|csc|sec|cot|arcsinh|arccosh|arctanh|arccsch|arcsech|arccoth|arcsin|arccos|arctan|arccsc|arcsec|arccot|log|lg|ln|exp|det|sqrt|pi|e|rad|deg|abs|inc|hirt|<=|>=|\\\"[^\\\"\\\\\\\\]*(\\\\\\\\.[^\\\"\\\\\\\\]*)*\\\"|[0-9.-]+|\\S")
	debug("using tokenizer regex:", tokenizerRegExp.String())
	multi := len(files) > 1
	for _, filename := range files {
		symbols = make(map[string]*Symbol)
		if multi {
			info(filename + ":")
		}
		ok, file := interpret(filename)
		if file != nil {
			file.Close()
		}
		if !ok {
			return
		}
		for k := range symbols {
			delete(symbols, k)
		}
	}
}

func interpret(filename string) (bool, *os.File) {
	file, e := os.Open(filename)
	if e != nil {
		err("could not open file:", filename)
		return false, file
	}
	scn := bufio.NewScanner(file)
	for line := 1; scn.Scan(); line++ {
		text := scn.Text()
		if !interpretLine(text, line) {
			return false, file
		}
	}
	file.Close()
	return true, nil
}

func interpretLine(data string, line int) bool {
	comment := strings.Index(data, "//")
	if comment != -1 {
		data = data[0:comment]
	}
	if len(strings.TrimSpace(data)) == 0 { // blank line
		return true
	}
	if isStringLiteral(data) {
		if data[0] == '"' {
			fmt.Println(unwrapStringLiteral(data))
		} else {
			fmt.Print(unwrapStringLiteral(data))
		}
		return true
	}
	iRad := strings.Index(data, "radians")
	iDeg := strings.Index(data, "degrees")
	if iRad != -1 {
		if strings.TrimSpace(data) != "radians" {
			errLine("the radians directive must be on a line by itself", line)
			return false
		}
		trigMode = Radians
		return true
	}
	if iDeg != -1 {
		if strings.TrimSpace(data) != "degrees" {
			errLine("the degrees directive must be on a line by itself", line)
			return false
		}
		trigMode = Degrees
		return true
	}
	tokens := tokenizerRegExp.FindAllString(data, -1)
	for i := 0; i < len(tokens)-1; i++ { // we love implicit multiplication!
		t1, t2 := tokens[i], tokens[i+1]
		if isOperator(t1) || isOperator(t2) {
			continue
		}
		if t1 == "(" || t2 == ")" || t1 == "," || t2 == "," {
			continue
		}
		_, st1e := symbols[t1]
		if isAlpha(t1) && t2 == "(" && (!st1e || symbols[t1].variant != constantVariant) {
			continue
		}
		tokens = append(tokens[:i+1], tokens[i:]...)
		tokens[i+1] = "*"
		i++
	}
	result := outputLineResult
	var fSymbol *Symbol = nil
	var cSymbol *Symbol = nil
	debug(strconv.Itoa(int(result)))
	if len(tokens) >= 2 && isAlpha(tokens[0]) && tokens[1] == "=" {
		result = constDefLineResult
		if len(tokens) == 2 {
			errLine("expression expected after equal (=) sign", line)
			return false
		}
		cName := tokens[0]
		tokens = tokens[2:]
		_, constExists := symbols[cName]
		if constExists {
			errLine("a constant or function is already declared with the name '"+cName+"'", line)
			return false
		}
		symbols[cName] = &Symbol{cName, constantVariant, []*Symbol{}, nil, Deque[string]{}}
		cSymbol = symbols[cName]
	}
	if len(tokens) >= 4 && isAlpha(tokens[0]) && tokens[1] == "(" {
		next := indexOfSA("(", tokens, 2)
		end := indexOfSA(")", tokens, 2)
		if (next == -1 || next > end) && end+1 < len(tokens) && tokens[end+1] == "=" {
			funcDef := tokens[:end+1]
			result = funcDefLineResult
			if end+2 >= len(tokens) {
				errLine("expression expected after equal (=) sign", line)
				return false
			}
			tokens = tokens[end+2:]
			defIt := 0
			fName := funcDef[defIt]
			_, funcExists := symbols[fName]
			if funcExists {
				errLine("a constant or function is already declared with the name '"+funcDef[defIt]+"'", line)
				return false
			}
			defIt++
			if funcDef[defIt] != "(" {
				errLine("expected left parenthesis to start function parameter list", line)
				return false
			}
			defIt++
			symbols[fName] = &Symbol{fName, functionVariant, []*Symbol{}, nil, Deque[string]{}}
			fSymbol = symbols[fName]
			for ; defIt < len(funcDef); defIt++ {
				current := funcDef[defIt]
				if !isAlpha(current) {
					errLine("a function parameter's name may only contain alphabetical characters", line)
					return false
				}
				foundShade, fsExists := symbols[current]
				var shaded *Symbol = nil
				if fsExists {
					shaded = foundShade
				}
				symbols[current] = &Symbol{current, functionParameterVariant, []*Symbol{}, shaded, Deque[string]{}}
				fSymbol.children = append(fSymbol.children, symbols[current])
				defIt++
				if defIt >= len(funcDef) || (funcDef[defIt] != "," && funcDef[defIt] != ")") {
					errLine("expected right parenthesis or comma", line)
				}
			}
		}
	}
	debug(strings.Join(tokens, ", "))
	output := Deque[string]{}
	ops := Stack[string]{}
	for i := 0; i < len(tokens); i++ {
		tp := &(tokens[i])
		fSymbol, fSymbolExists := symbols[*tp]
		if fSymbolExists && fSymbol.variant != functionVariant {
			fSymbolExists = false
		}
		trans, transExists := transcendentals[*tp]
		if isNumeric(*tp) {
			output.push(tp)
		} else if fSymbolExists || (transExists && trans.variant == functionVariant) {
			ops.push(tp)
		} else if isOperator(*tp) {
			for !ops.empty() && *(ops.peek()) != "(" && (operators[*(ops.peek())] > operators[*tp] || operators[*(ops.peek())] == operators[*tp] && *tp != "^") {
				output.push(ops.pop())
			}
			ops.push(tp)
		} else if *tp == "(" {
			ops.push(tp)
		} else if *tp == ")" {
			for !ops.empty() && *(ops.peek()) != "(" {
				output.push(ops.pop())
			}
			if ops.empty() || *(ops.peek()) != "(" {
				err("expected left parenthesis")
				return false
			}
			ops.pop()
			if !ops.empty() && isAlpha(*(ops.peek())) {
				output.push(ops.pop())
			}
		} else if *tp != "," {
			output.push(tp)
		}
		debug(output.string())
	}
	for !ops.empty() {
		if *(ops.peek()) == "(" {
			err("unexpected left parenthesis")
			return false
		}
		output.push(ops.pop())
		debug(output.string())
	}
	debug(output.string())

	if result == funcDefLineResult {
		fSymbol.data = output
	} else {
		value, no := resolveExpression(&output, nil, nil, line)
		if result == constDefLineResult {
			cSymbol.data = Deque[string]{}
			cSymbol.data.push(&value)
		} else {
			result, _ := strconv.ParseFloat(value, 64)
			if result > 0 && result <= 0.00001 {
				result = 0
			}
			if !no {
				if !math.IsNaN(result) && !math.IsInf(result, 0) {
					fmt.Println(result)
				} else {
					fmt.Println("undefined")
				}
			}
		}
	}
	// clean up
	if fSymbol != nil {
		// delete all children (function parameters) from symbol table or replace with their shaded versions
		for _, child := range fSymbol.children {
			if child.shaded != nil {
				symbols[child.name] = child.shaded
				child.shaded = nil
			} else {
				delete(symbols, child.name)
			}
		}
	}
	for k := range symbols {
		debug(k, symbols[k].string())
	}
	return true
}

func resolveExpression(data *Deque[string], functionChildren *[]*Symbol, operations *Stack[string], line int) (string, bool) {
	valueTable := make(map[string]string)
	if functionChildren != nil && operations != nil {
		for i := len(*functionChildren) - 1; i >= 0; i-- {
			valueTable[(*functionChildren)[i].name] = *(operations.pop())
		}
	}
	localOperations := Stack[string]{}
	noOutput := false
	for current := data.shift(); current != nil; current = data.shift() {
		symbol, fSymbolExists := symbols[*current]
		cSymbolExists := fSymbolExists
		if fSymbolExists && symbol.variant != functionVariant {
			fSymbolExists = false
		}
		value, valueExists := valueTable[*current]
		trans, transExists := transcendentals[*current]
		if fSymbolExists {
			data := symbol.data
			result, nol := resolveExpression(&data, &(symbol.children), &localOperations, line)
			if nol {
				noOutput = true
			}
			localOperations.push(&result)
		} else if transExists {
			args := make([]string, 0, trans.argCount)
			for i := 0; i < trans.argCount; i++ {
				arg := localOperations.pop()
				if arg == nil {
					errLine("function "+*current+" expects "+strconv.Itoa(trans.argCount)+" argument(s), got "+strconv.Itoa(i), line)
					panic("-1")
				}
				args = append(args, *arg)
			}
			result := trans.operation(args...)
			if result == "" {
				result = "0"
				noOutput = true
			}
			localOperations.push(&result)
		} else if valueExists {
			localOperations.push(&value)
		} else if cSymbolExists {
			localOperations.push(symbol.data.front())
		} else if isNumeric(*current) {
			localOperations.push(current)
		} else if isStringLiteral(*current) {
			value := unwrapStringLiteral(*current)
			localOperations.push(&value)
		} else if isOperator(*current) {
			rhs := localOperations.pop()
			lhs := localOperations.pop()
			if lhs == nil || rhs == nil {
				errLine("'"+*current+"' operator expects 2 operands", line)
				panic("-1")
			}
			lhsv, lhsve := strconv.ParseFloat(*lhs, 64)
			rhsv, rhsve := strconv.ParseFloat(*rhs, 64)
			if lhsve != nil || rhsve != nil {
				errLine("'"+*current+"' operator expects 2 operands", line)
				panic("-1")
			}
			value := ""
			switch *current {
			case "+":
				value = strconv.FormatFloat(lhsv+rhsv, 'E', -1, 64)
			case "-":
				value = strconv.FormatFloat(lhsv-rhsv, 'E', -1, 64)
			case "*":
				value = strconv.FormatFloat(lhsv*rhsv, 'E', -1, 64)
			case "/":
				value = strconv.FormatFloat(lhsv/rhsv, 'E', -1, 64)
			case "^":
				value = strconv.FormatFloat(math.Pow(lhsv, rhsv), 'E', -1, 64)
			default:
				{
					errLine("'"+*current+"' operator has not been implemented yet", line)
					panic("-1")
				}
			}
			localOperations.push(&value)
		} else {
			errLine("unknown symbol encountered: "+*current, line)
			panic("-1")
		}
	}
	if localOperations.empty() {
		errLine("there was an issue while resolving an expression", line)
		panic("-1")
	}
	return *(localOperations.pop()), noOutput
}
