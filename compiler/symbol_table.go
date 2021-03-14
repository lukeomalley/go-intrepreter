package compiler

// SymbolScope stores the scope of the symbol
type SymbolScope string

// Enum of possible symbol scopes
const (
	GlobalScope  SymbolScope = "GLOBAL"
	LocalScope   SymbolScope = "LOCAL"
	BuiltinScope SymbolScope = "BUILTIN"
	FreeScope    SymbolScope = "FREE"
)

// Symbol represents an identifier in the users program
type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

// SymbolTable stores a reference from the original name to the symbol
type SymbolTable struct {
	Outer          *SymbolTable
	store          map[string]Symbol
	numDefinitions int
	FreeSymbols    []Symbol
}

// NewSymbolTable constructs a symbol table
func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	free := []Symbol{}
	return &SymbolTable{store: s, FreeSymbols: free}
}

// NewEnclosedSymbolTable constructs a symbol table for an enclosed scope
func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.Outer = outer
	return s
}

// Define adds a symbol to the symbol table
func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions}
	if s.Outer != nil {
		symbol.Scope = LocalScope
	} else {
		symbol.Scope = GlobalScope
	}

	s.store[name] = symbol
	s.numDefinitions++
	return symbol
}

// Resolve looks up a symbol in the symbol table
func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]

	if !ok && s.Outer != nil {
		obj, ok := s.Outer.Resolve(name)
		if !ok {
			return obj, ok
		}

		if obj.Scope == GlobalScope || obj.Scope == BuiltinScope {
			return obj, ok
		}

		free := s.defineFree(obj)
		return free, true
	}

	return obj, ok
}

func (s *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	symbol := Symbol{Name: name, Index: index, Scope: BuiltinScope}
	s.store[name] = symbol
	return symbol
}

func (s *SymbolTable) defineFree(original Symbol) Symbol {
	s.FreeSymbols = append(s.FreeSymbols, original)

	symbol := Symbol{Name: original.Name, Index: len(s.FreeSymbols) - 1}
	symbol.Scope = FreeScope // TODO: Why cant this be done in one line?

	s.store[original.Name] = symbol
	return symbol
}
