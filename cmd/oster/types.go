package main

import (
	"fmt"
	"strings"
)

type traceContext struct {
	binaryName   string
	functionName string
	Arguments    []argument
}

type argument struct {
	CType          string
	goType         goType
	StartingOffset int
	VariableName   string
	PrintfFormat   string
}

type goType int

const (
	INVALID        = 0
	INT     goType = iota
	INT8
	INT16
	INT32
	INT64
	UINT
	UINT8
	UINT16
	UINT32
	UINT64
	FLOAT32
	FLOAT64
	BOOL
	STRING
	BYTE
	//TODO:
	STRUCT
	POINTER
)

var goTypeToSizeInBytes = map[goType]int{
	INT:     8,
	INT8:    1,
	INT16:   2,
	INT32:   4,
	INT64:   8,
	UINT:    8,
	UINT8:   1,
	UINT16:  2,
	UINT32:  4,
	UINT64:  8,
	FLOAT32: 4,
	FLOAT64: 8,
	BOOL:    1,
	BYTE:    1,
	STRING:  8,
	//TODO:
	STRUCT:  8,
	POINTER: 8,
}

var goToCType = map[goType]string{
	INT:     "long",
	INT8:    "char",
	INT16:   "short",
	INT32:   "int",
	INT64:   "long",
	UINT:    "long",
	UINT8:   "char",
	UINT16:  "short",
	UINT32:  "int",
	UINT64:  "long",
	FLOAT32: "float",
	FLOAT64: "double",
	BOOL:    "char",
	BYTE:    "char",
	STRING:  "char *",
	//TODO:
	STRUCT:  "void *",
	POINTER: "void *",
}

func stringfFormat(t goType) string {
	switch t {
	case INT8, INT16, INT32, UINT8, UINT16, UINT32:
		return "%d"
	case INT, UINT, INT64, UINT64:
		return "%ld"
	case FLOAT32, FLOAT64:
		return "%e"
	case BOOL:
		return "%t"
	case STRING:
		return "%s"
	case BYTE:
		return "%c"
	//TODO:
	case STRUCT, POINTER:
		return "0x%x"
	default:
		return "%v"
	}
}

var stringToGoType = map[string]goType{
	"INT":     INT,
	"INT8":    INT8,
	"INT16":   INT16,
	"INT32":   INT32,
	"INT64":   INT64,
	"UINT":    UINT,
	"UINT8":   UINT8,
	"UINT16":  UINT16,
	"UINT32":  UINT32,
	"UINT64":  UINT64,
	"FLOAT32": FLOAT32,
	"FLOAT64": FLOAT64,
	"BOOL":    BOOL,
	"STRING":  STRING,
	"BYTE":    BYTE,
	//TODO:
	"STRUCT":  STRUCT,
	"POINTER": POINTER,
}

var goTypeToString = map[goType]string{
	INT:     "INT",
	INT8:    "INT8",
	INT16:   "INT16",
	INT32:   "INT32",
	INT64:   "INT64",
	UINT:    "UINT",
	UINT8:   "UINT8",
	UINT16:  "UINT16",
	UINT32:  "UINT32",
	UINT64:  "UINT64",
	FLOAT32: "FLOAT32",
	FLOAT64: "FLOAT64",
	BOOL:    "BOOL",
	STRING:  "STRING",
	BYTE:    "BYTE",
	//TODO:
	STRUCT:  "STRUCT",
	POINTER: "POINTER",
}

func listAvailableTypes() {
	for k, _ := range stringToGoType {
		fmt.Println(k)
	}
}

type stack []byte

var invalidChars = "+&%$#@!<>/?\";:[]{}=-`~" //fixme: this isn't exhaustive, doesn't take into account digits as first char

func (s *stack) push(v byte) bool {

	if v == ' ' {
		return true
	}

	*s = append(*s, v)
	return true
}

func (s *stack) clear() {
	*s = []byte{}
}

func (s *stack) string() string {
	return string(*s)
}

// parseFunctionAndArgumentTypes populates the traceContext based on the function and argument types
// of the form 'func_name(type1, type2)'.
func parseFunctionAndArgumentTypes(context *traceContext, funcAndArgs string) error {

	parseStack := &stack{}

	for i := range funcAndArgs {

		if strings.ContainsAny(string(funcAndArgs[i]), invalidChars) {
			return fmt.Errorf("encountered invalid char: %s", string(funcAndArgs[i]))
		}

		if funcAndArgs[i] == '(' {
			context.functionName = parseStack.string()
			invalidChars += string('(')
			parseStack.clear()
			continue
		}

		if funcAndArgs[i] == ',' {
			goType := stringToGoType[strings.ToUpper(parseStack.string())]
			if goType == INVALID {
				return fmt.Errorf("invalid go type: %s", parseStack.string())
			}

			newArg := argument{
				goType:       goType,
				PrintfFormat: stringfFormat(goType),
				CType:        goToCType[goType],
				VariableName: fmt.Sprintf("argument%d", i), //variable name doesn't actually matter, just needs to be unique
			}

			context.Arguments = append(context.Arguments, newArg)
			parseStack.clear()
			continue
		}

		if funcAndArgs[i] == ')' {
			goType := stringToGoType[strings.ToUpper(parseStack.string())]
			if goType == INVALID {
				return fmt.Errorf("invalid go type: %s", parseStack.string())
			}
			newArg := argument{
				goType:       goType,
				PrintfFormat: stringfFormat(goType),
				CType:        goToCType[goType],
				VariableName: fmt.Sprintf("argument%d", i),
			}

			context.Arguments = append(context.Arguments, newArg)
			return nil
		}

		parseStack.push(funcAndArgs[i])
	}

	return nil
}
