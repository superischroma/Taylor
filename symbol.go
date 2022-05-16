package main

import (
	"strconv"
	"strings"
)

type SymbolVariant int
type ScopeVariant int

const (
	constantVariant          SymbolVariant = 0
	functionVariant          SymbolVariant = 1
	functionParameterVariant SymbolVariant = 2

	ifVariant   ScopeVariant = 0
	loopVariant ScopeVariant = 1
)

type Symbol struct {
	name     string
	variant  SymbolVariant
	children []*Symbol
	shaded   *Symbol
	data     Deque[string]
}

type Scope struct {
	scope   *Scope
	variant ScopeVariant
	data    []string
}

func (symbol *Symbol) string() string {
	builder := strings.Builder{}
	builder.WriteString("Symbol{name=")
	builder.WriteString(symbol.name)
	builder.WriteString(", variant=")
	builder.WriteString(strconv.Itoa(int(symbol.variant)))
	builder.WriteString(", children=[")
	for i, child := range symbol.children {
		if i != 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(child.string())
	}
	builder.WriteString("], shaded=")
	if symbol.shaded != nil {
		builder.WriteString(symbol.shaded.string())
	} else {
		builder.WriteString("<nil>")
	}
	builder.WriteRune('}')
	return builder.String()
}
