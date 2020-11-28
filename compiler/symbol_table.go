package compiler

// SymbolScope stores the scope of the symbol
type SymbolScope string

// Enum of possible symbol scopes
const (
	GlobalScope SymbolScope = "GLOBAL"
)

// Symbol represents an identifier in the users program
type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

// SymbolTable stores a reference from the original name to the symbol
type SymbolTable struct {
	store          map[string]Symbol
	numDefinitions int
}

// NewSymbolTable constructs a symbol table
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store:          make(map[string]Symbol),
		numDefinitions: 0,
	}
}

// Define adds a symbol to the symbol table
func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions, Scope: GlobalScope}
	s.store[name] = symbol
	s.numDefinitions++
	return symbol
}

// Resolve looks up a symbol in the symbol table
func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]
	return obj, ok
}
