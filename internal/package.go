package internal

import "strings"

type Package struct {
	Name string
	Doc  string

	Variables      []Variable
	VariableBlocks []VariableBlock
	Constants      []Variable
	ConstantBlocks []VariableBlock

	Functions []Function

	Types      []Type
	Structs    []Struct
	Interfaces []Interface
}

type Function struct {
	Name       string
	Doc        string
	Definition string
}

func (i *Function) addToDocs(docs string) {
	i.Doc += strings.TrimLeft(docs, " ") + "\n"
}

type Variable struct {
	Name       string
	Doc        string
	Definition string
	Value      string
	Type       string
}

func (i *Variable) addToDocs(docs string) {
	i.Doc += strings.TrimLeft(docs, " ") + "\n"
}

type VariableBlock struct {
	Variables []Variable
	Doc       string
}

func (i *VariableBlock) addToDocs(docs string) {
	i.Doc += strings.TrimLeft(docs, " ") + "\n"
}

type Type struct {
	Doc        string
	Name       string
	Definition string
	Functions  []Function
}

func (i *Type) addToDocs(docs string) {
	i.Doc += strings.TrimLeft(docs, " ") + "\n"
}

type Struct struct {
	Doc        string
	Name       string
	Definition string
	Functions  []Function
}

func (i *Struct) addToDocs(docs string) {
	i.Doc += strings.TrimLeft(docs, " ") + "\n"
}

type Interface struct {
	Doc        string
	Name       string
	Definition string
	Values     []Variable
}

func (i *Interface) addToDocs(docs string) {
	i.Doc += strings.TrimLeft(docs, " ") + "\n"
}

type documentable interface {
	addToDocs(string)
}

func (p Package) getLastVariable() *Variable {
	if len(p.Variables) > 0 {
		return &p.Variables[len(p.Variables)-1]
	}

	return &Variable{}
}

func (p Package) getLastVariableBlock() *VariableBlock {
	if len(p.VariableBlocks) > 0 {
		return &p.VariableBlocks[len(p.VariableBlocks)-1]
	}

	return &VariableBlock{}
}

func (p Package) getLastConstant() *Variable {
	if len(p.Constants) > 0 {
		return &p.Constants[len(p.Constants)-1]
	}

	return &Variable{}
}

func (p Package) getLastConstantBlock() *VariableBlock {
	if len(p.ConstantBlocks) > 0 {
		return &p.ConstantBlocks[len(p.ConstantBlocks)-1]
	}

	return &VariableBlock{}
}

func (p Package) getLastFunction() *Function {
	if len(p.Functions) > 0 {
		return &p.Functions[len(p.Functions)-1]
	}

	return &Function{}
}

func (p Package) getLastType() *Type {
	if len(p.Types) > 0 {
		return &p.Types[len(p.Types)-1]
	}

	return &Type{}
}

func (p Package) getLastStruct() *Struct {
	if len(p.Structs) > 0 {
		return &p.Structs[len(p.Structs)-1]
	}

	return &Struct{}
}

func (p Package) getLastInterface() *Interface {
	if len(p.Interfaces) > 0 {
		return &p.Interfaces[len(p.Interfaces)-1]
	}

	return &Interface{}
}
