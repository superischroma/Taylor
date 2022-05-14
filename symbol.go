package main

type SymbolVariant int

const (
	constantVariant SymbolVariant = 0
	functionVariant SymbolVariant = 1
)

type Symbol struct {
	name    string
	variant SymbolVariant
}
