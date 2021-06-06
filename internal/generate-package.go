package internal

func GenerateTestPackage() Package {
	return Package{
		Name: "experimenting",
		Doc:  "Package experimenting is an experimenting package.\nThis is the package doc.",
		Variables: []Variable{
			{Name: "SimpleEmptyString", Doc: "SimpleEmptyString doc.", Definition: `var SimpleEmptyString string`, Value: "", Type: "string"},
			{Name: "SimpleEmptyInt", Doc: "SimpleEmptyInt doc.", Definition: `var SimpleEmptyInt int`, Value: "", Type: "int"},
			{Name: "ThisIsASingleStringVariable", Doc: "ThisIsASingleStringVariable doc.", Definition: `var ThisIsASingleStringVariable = "Hello, World!"`, Value: `"Hello, World!"`, Type: ""},
			{Name: "ThisIsASingleIntVariable", Doc: "ThisIsASingleIntVariable doc.", Definition: `var ThisIsASingleIntVariable = 1337`, Value: "1337", Type: ""},
			{Name: "SS1Var", Doc: "SS1Var doc.", Definition: `var SS1Var SimpleStruct1`, Value: "", Type: "SimpleStruct1"},
			{Name: "SS2Var", Doc: "SS2Var doc.", Definition: `var SS2Var SimpleStruct2`, Value: "", Type: "SimpleStruct2"},
		},
		VariableBlocks: []VariableBlock{
			{
				Variables: []Variable{
					{Name: "ThisIsInsideTheFirstVariableBlock1", Doc: "ThisIsInsideTheFirstVariableBlock1 docs.", Definition: `ThisIsInsideTheFirstVariableBlock1 = ""`, Value: "", Type: ""},
					{Name: "ThisIsInsideTheFirstVariableBlock2", Doc: "ThisIsInsideTheFirstVariableBlock2 docs.", Definition: `ThisIsInsideTheFirstVariableBlock2 = ""`, Value: "", Type: ""},
					{Name: "ThisIsInsideTheFirstVariableBlock3", Doc: "ThisIsInsideTheFirstVariableBlock3 docs.", Definition: `ThisIsInsideTheFirstVariableBlock3 = ""`, Value: "", Type: ""},
					{Name: "ThisIsInsideTheFirstVariableBlock4", Doc: "ThisIsInsideTheFirstVariableBlock4 docs.", Definition: `ThisIsInsideTheFirstVariableBlock4 = ""`, Value: "", Type: ""},
					{Name: "ThisIsInsideTheFirstVariableBlock5", Doc: "ThisIsInsideTheFirstVariableBlock5 docs.", Definition: `ThisIsInsideTheFirstVariableBlock5 = ""`, Value: "", Type: ""},
				},
				Doc: "docs.",
			},
			{
				Variables: []Variable{
					{Name: "ThisIsInsideTheSecondVariableBlock1", Doc: "ThisIsInsideTheSecondVariableBlock1 docs.", Definition: `ThisIsInsideTheSecondVariableBlock1 = ""`, Value: "", Type: ""},
					{Name: "ThisIsInsideTheSecondVariableBlock2", Doc: "ThisIsInsideTheSecondVariableBlock2 docs.", Definition: `ThisIsInsideTheSecondVariableBlock2 = ""`, Value: "", Type: ""},
					{Name: "ThisIsInsideTheSecondVariableBlock3", Doc: "ThisIsInsideTheSecondVariableBlock3 docs.", Definition: `ThisIsInsideTheSecondVariableBlock3 = ""`, Value: "", Type: ""},
					{Name: "ThisIsInsideTheSecondVariableBlock4", Doc: "ThisIsInsideTheSecondVariableBlock4 docs.", Definition: `ThisIsInsideTheSecondVariableBlock4 = ""`, Value: "", Type: ""},
					{Name: "ThisIsInsideTheSecondVariableBlock5", Doc: "ThisIsInsideTheSecondVariableBlock5 docs.", Definition: `ThisIsInsideTheSecondVariableBlock5 = ""`, Value: "", Type: ""},
				},
				Doc: "docs.",
			},
		},
		Constants: []Variable{
			{Name: "ThisIsASingleStringConst", Doc: "ThisIsASingleStringConst docs.", Definition: `const ThisIsASingleStringConst = "Hello, World!"`, Value: `"Hello, World!""`, Type: ""},
			{Name: "ThisIsASingleIntConst", Doc: "ThisIsASingleIntConst docs.", Definition: `const ThisIsASingleIntConst = 1337`, Value: `1337`, Type: ""},
		},
		ConstantBlocks: []VariableBlock{
			{
				Variables: []Variable{
					{Name: "ThisIsInsideTheFirstConstantBlock1", Doc: "ThisIsInsideTheFirstConstantBlock1 docs.", Definition: `ThisIsInsideTheFirstConstantBlock1 = ""`, Value: "", Type: ""},
					{Name: "ThisIsInsideTheFirstConstantBlock2", Doc: "ThisIsInsideTheFirstConstantBlock2 docs.", Definition: `ThisIsInsideTheFirstConstantBlock2 = ""`, Value: "", Type: ""},
					{Name: "ThisIsInsideTheFirstConstantBlock3", Doc: "ThisIsInsideTheFirstConstantBlock3 docs.", Definition: `ThisIsInsideTheFirstConstantBlock3 = ""`, Value: "", Type: ""},
					{Name: "ThisIsInsideTheFirstConstantBlock4", Doc: "ThisIsInsideTheFirstConstantBlock4 docs.", Definition: `ThisIsInsideTheFirstConstantBlock4 = ""`, Value: "", Type: ""},
					{Name: "ThisIsInsideTheFirstConstantBlock5", Doc: "ThisIsInsideTheFirstConstantBlock5 docs.", Definition: `ThisIsInsideTheFirstConstantBlock5 = ""`, Value: "", Type: ""},
				},
				Doc: "",
			},
			{
				Variables: []Variable{
					{Name: "ThisIsInsideTheSecondConstantBlock1", Doc: "ThisIsInsideTheSecondConstantBlock1 docs.", Definition: `ThisIsInsideTheSecondConstantBlock1 = ""`, Value: "", Type: ""},
					{Name: "ThisIsInsideTheSecondConstantBlock2", Doc: "ThisIsInsideTheSecondConstantBlock2 docs.", Definition: `ThisIsInsideTheSecondConstantBlock2 = ""`, Value: "", Type: ""},
					{Name: "ThisIsInsideTheSecondConstantBlock3", Doc: "ThisIsInsideTheSecondConstantBlock3 docs.", Definition: `ThisIsInsideTheSecondConstantBlock3 = ""`, Value: "", Type: ""},
					{Name: "ThisIsInsideTheSecondConstantBlock4", Doc: "ThisIsInsideTheSecondConstantBlock4 docs.", Definition: `ThisIsInsideTheSecondConstantBlock4 = ""`, Value: "", Type: ""},
					{Name: "ThisIsInsideTheSecondConstantBlock5", Doc: "ThisIsInsideTheSecondConstantBlock5 docs.", Definition: `ThisIsInsideTheSecondConstantBlock5 = ""`, Value: "", Type: ""},
				},
				Doc: "docs.",
			},
		},
		Functions: []Function{
			{Name: "DoesNothing", Doc: "DoesNothing docs.", Definition: "func DoesNothing()"},
			{Name: "ReturnsTwoThings", Doc: "ReturnsTwoThings docs.", Definition: "func ReturnsTwoThings() (string, err)"},
		},
		Types: []Type{
			{Doc: "ST docs.", Name: "ST", Definition: `type ST string`, Functions: []Function{
				{Name: "DoesNothing", Doc: "DoesNothing docs.", Definition: "func (t *ST) DoesNothing()"},
				{Name: "ReturnsTwoThings", Doc: "ReturnsTwoThings docs.", Definition: "func (t *ST) ReturnsTwoThings() (string, err)"},
			}},
		},
		Structs: []Struct{
			{
				Doc:  "SS docs.",
				Name: "SS",
				Definition: `type SS struct {
	Name string
	Age int
}`,
				Functions: []Function{{Name: "DoesNothing", Doc: "DoesNothing docs.", Definition: "func (t *SS) DoesNothing()"},
					{Name: "ReturnsTwoThings", Doc: "ReturnsTwoThings docs.", Definition: "func (t *SS) ReturnsTwoThings() (string, err)"}},
			},
		},
		Interfaces: []Interface{
			{
				Doc:  "SI docs.",
				Name: "SI",
				Definition: `type SI interface {
	Run(text string)
	Get() string
}`,
				Values: []Variable{
					{Name: "Run", Doc: "", Definition: "Run(text string)", Value: "", Type: ""},
					{Name: "Get", Doc: "", Definition: "Get() string", Value: "", Type: ""},
				},
			},
		},
	}
}

const TestGodoc = `package experimenting // import "github.com/MarvinJWendt/gomark/experimenting"

Package experimenting is an experimenting package. This is the package doc.

CONSTANTS

const (
	// ThisIsInsideTheFirstConstantBlock1 docs.
	ThisIsInsideTheFirstConstantBlock1 = ""
	// ThisIsInsideTheFirstConstantBlock2 docs.
	ThisIsInsideTheFirstConstantBlock2 = ""
	// ThisIsInsideTheFirstConstantBlock3 docs.
	ThisIsInsideTheFirstConstantBlock3 = ""
	// ThisIsInsideTheFirstConstantBlock4 docs.
	ThisIsInsideTheFirstConstantBlock4 = ""
	// ThisIsInsideTheFirstConstantBlock5 docs.
	ThisIsInsideTheFirstConstantBlock5 = ""
)
    docs.

const (
	// ThisIsInsideTheSecondConstantBlock1 docs.
	ThisIsInsideTheSecondConstantBlock1 = ""
	// ThisIsInsideTheSecondConstantBlock2 docs.
	ThisIsInsideTheSecondConstantBlock2 = ""
	// ThisIsInsideTheSecondConstantBlock3 docs.
	ThisIsInsideTheSecondConstantBlock3 = ""
	// ThisIsInsideTheSecondConstantBlock4 docs.
	ThisIsInsideTheSecondConstantBlock4 = ""
	// ThisIsInsideTheSecondConstantBlock5 docs.
	ThisIsInsideTheSecondConstantBlock5 = ""
)
    docs.

const ThisIsASingleIntConst = 1337
    ThisIsASingleIntConst docs.

const ThisIsASingleStringConst = "Hello, World!"
    ThisIsASingleStringConst docs.


VARIABLES

var (
	// ThisIsInsideTheFirstVariableBlock1 docs.
	ThisIsInsideTheFirstVariableBlock1 = ""
	// ThisIsInsideTheFirstVariableBlock2 docs.
	ThisIsInsideTheFirstVariableBlock2 = ""
	// ThisIsInsideTheFirstVariableBlock3 docs.
	ThisIsInsideTheFirstVariableBlock3 = ""
	// ThisIsInsideTheFirstVariableBlock4 docs.
	ThisIsInsideTheFirstVariableBlock4 = ""
	// ThisIsInsideTheFirstVariableBlock5 docs.
	ThisIsInsideTheFirstVariableBlock5 = ""
)
    docs.

var (
	// ThisIsInsideTheSecondVariableBlock1 docs.
	ThisIsInsideTheSecondVariableBlock1 = ""
	// ThisIsInsideTheSecondVariableBlock2 docs.
	ThisIsInsideTheSecondVariableBlock2 = ""
	// ThisIsInsideTheSecondVariableBlock3 docs.
	ThisIsInsideTheSecondVariableBlock3 = ""
	// ThisIsInsideTheSecondVariableBlock4 docs.
	ThisIsInsideTheSecondVariableBlock4 = ""
	// ThisIsInsideTheSecondVariableBlock5 docs.
	ThisIsInsideTheSecondVariableBlock5 = ""
)
    docs.

var SimpleEmptyInt int
    SimpleEmptyInt doc.

var SimpleEmptyString string
    SimpleEmptyString doc.

var ThisIsASingleIntVariable = 1337
    ThisIsASingleIntVariable docs.

var ThisIsASingleStringVariable = "Hello, World!"
    ThisIsASingleStringVariable docs.


FUNCTIONS

func AcceptsAString(input string)
    AcceptsAString docs.

func DoesNothing()
    DoesNothing docs.

func ReturnPointerMap() map[*string]*error
    ReturnPointerMap docs.

func ReturnsAString() string
    ReturnsAString docs.

func ReturnsPointerInt() *int
    ReturnsPointerInt docs.

func TakesAndReturnsAString(input string) string
    TakesAndReturnsAString docs.


TYPES

type SS1Slice []SimpleStruct1
    SS1Slice docs.

type SimpleInterface1 interface {
	Run()
	Get() string
}
    SimpleInterface1 docs.

type SimpleInterface2 interface {
	Set(string)
	Get() string
	Write() int
}
    SimpleInterface2 docs.

type SimpleStruct1 struct {
	StringField string
	IntField    int
}
    SimpleStruct1 docs.

var SS1Var SimpleStruct1
    SS1Var docs.

func (s SimpleStruct1) SS1F1(in string) string
    SS1F1 docs.

func (s SimpleStruct1) SS1F2(in string) int
    SS1F2 docs.

type SimpleStruct2 struct {
	StringField string
	IntField    int
}
    SimpleStruct2 docs.

var SS2Var SimpleStruct2
    SS2Var docs.

func (s SimpleStruct2) SS2F1(in string) string
    SS2F1 docs.

func (s SimpleStruct2) SS2F2(in string) int
    SS2F2 docs.

type TypeString string
    TypeString docs.
`
