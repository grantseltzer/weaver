package main

import (
	"debug/dwarf"
	"debug/elf"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// TraceTarget represents a function/method to attach a uprobe+ebpf to
type TraceTarget struct {
	Name       string
	Parameters []Parameter
	Returns    []Parameter
}

// Parameter represents a parameter
type Parameter struct {
	Name         string
	TypeString   string
	GoType       goType
	CType        string
	StartOffset  int
	PrintfFormat string
	TypeSize     int
	ArrayLength  int // 0 if not an array
}

func getDwarfData(path string) (*dwarf.Data, error) {
	elfFile, err := elf.Open(path)
	if err != nil {
		return nil, err
	}

	dwarfData, err := elfFile.DWARF()
	if err != nil {
		return nil, err
	}

	return dwarfData, nil
}

func enrichTargets(targets []TraceTarget) error {
	for i := range targets {
		for n := range targets[i].Parameters {
			err := enrichParameter(&targets[i].Parameters[n])
			if err != nil {
				return err
			}
		}

		for m := range targets[i].Returns {
			err := enrichParameter(&targets[i].Returns[m])
			if err != nil {
				return err
			}
		}

		err := getStackOffsets(&targets[i])
		if err != nil {
			return err
		}
	}

	return nil
}

// parseDwarfData takes DWARF data and returns a slice
// of TraceTargets for weaver to attach uprobes/ebpf to.
func parseDwarfData(data *dwarf.Data) ([]TraceTarget, error) {

	linearReader := data.Reader()
	typeReader := data.Reader()

	var targets []TraceTarget

	var targetBeingRead *TraceTarget = nil

entryReadLoop:
	for {
		entry, err := linearReader.Next()
		if err == io.EOF || entry == nil {
			break entryReadLoop
		}
		if err != nil {
			return []TraceTarget{}, err
		}

		if targetBeingRead != nil {
			// currently reading in the parameters of a function symbol

			// Null entry is used to end function's list of parameters/variables
			if entryIsNull(entry) {
				targets = append(targets, *targetBeingRead)
				targetBeingRead = nil
				continue entryReadLoop
			}

			if entry.Tag != dwarf.TagFormalParameter {
				// Don't care about variables in func body, only parameters and returns
				continue entryReadLoop
			}

			newParam := Parameter{}
			isReturn := false

			// Get this parameter's name and type
			for i := range entry.Field {

				if entry.Field[i].Attr == dwarf.AttrName {
					newParam.Name = entry.Field[i].Val.(string)
				}

				if entry.Field[i].Attr == dwarf.AttrVarParam {
					if entry.Field[i].Val.(bool) == true {
						isReturn = true
					}
				}

				if entry.Field[i].Attr == dwarf.AttrType {
					typeReader.Seek(entry.Field[i].Val.(dwarf.Offset))
					typeEntry, err := typeReader.Next()
					if err != nil {
						return []TraceTarget{}, err
					}

					for i := range typeEntry.Field {
						if typeEntry.Field[i].Attr == dwarf.AttrName {
							newParam.TypeString = typeEntry.Field[i].Val.(string)
						}
					}
				}
			}

			if isReturn {
				targetBeingRead.Returns = append(targetBeingRead.Returns, newParam)
			} else {
				targetBeingRead.Parameters = append(targetBeingRead.Parameters, newParam)
			}
		}

		// debug entry is a function/method symbol
		if entry.Tag == dwarf.TagSubprogram {

			targetBeingRead = &TraceTarget{}

			// collect the symbols name by finding it in the entry fields
			for i := range entry.Field {
				if entry.Field[i].Attr == dwarf.AttrName {
					targetBeingRead.Name = entry.Field[i].Val.(string)
				}
			}
		}
	}
	return targets, nil
}

func entryIsNull(e *dwarf.Entry) bool {
	return e.Children == false &&
		len(e.Field) == 0 &&
		e.Offset == 0 &&
		e.Tag == dwarf.Tag(0)
}

func enrichParameter(param *Parameter) error {
	if strings.Contains(param.TypeString, "[") {
		size, gotype, err := parseArrayString(param.TypeString)
		if err != nil {
			return err
		}
		param.TypeSize = size
		param.GoType = gotype
	} else {
		goType := stringToGoType[param.TypeString]
		if goType == INVALID {
			return fmt.Errorf("invalid go type: %s", param.TypeString)
		}
		param.GoType = goType
	}

	param.CType = goToCType[param.GoType]

	return nil
}

// Determining stack sizes:
//
// - Look at size of largest data type that's being passed, that sets the window size
// - Each element added is limited by whether or not it will fit into that window
// - If it would go over a limit window then pad until back at 0, add it, then continue
func getStackOffsets(target *TraceTarget) error {

	var windowSize = 0

	for _, t := range target.Parameters {
		size := goTypeToSizeInBytes[t.GoType]
		if size > windowSize {
			windowSize = size
		}
	}

	currentIndex := 8
	bytesInCurrentWindow := 0

	for i := range target.Parameters {

		size := goTypeToSizeInBytes[target.Parameters[i].GoType]
		target.Parameters[i].TypeSize = size

		if size+bytesInCurrentWindow > windowSize {
			// Doesn't fit, move index ahead for padding, clear current window
			currentIndex += windowSize - bytesInCurrentWindow
			bytesInCurrentWindow = 0
		}

		target.Parameters[i].StartOffset = currentIndex

		if target.Parameters[i].ArrayLength > 0 {
			if target.Parameters[i].GoType == STRING {
				size = 16
			}
			currentIndex += size * target.Parameters[i].ArrayLength
			bytesInCurrentWindow += (size * target.Parameters[i].ArrayLength) % windowSize
			continue
		}

		currentIndex += size
		bytesInCurrentWindow += size

		//XXX: In go strings take up 16 bytes on the stack, 8 for the pointer and 8 for length
		if target.Parameters[i].GoType == STRING {
			currentIndex += 8
		}
	}
	return nil
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
	RUNE
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
	RUNE:    4,
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
	RUNE:    "int",
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
	case RUNE:
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
	"RUNE":    RUNE,
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
	RUNE:    "RUNE",
	//TODO:
	STRUCT:  "STRUCT",
	POINTER: "POINTER",
}

var supportedTypes = []string{
	"INT",
	"INT8",
	"INT16",
	"INT32",
	"INT64",
	"UINT",
	"UINT8",
	"UINT16",
	"UINT32",
	"UINT64",
	"FLOAT32",
	"FLOAT64",
	"BOOL",
	"STRING",
	"BYTE",
	"RUNE",
}

func parseArrayString(s string) (int, goType, error) {
	subs := strings.Split(s, "[")
	if len(subs) != 2 && subs[0] != "" {
		return -1, INVALID, errors.New("malformed array parameter")
	}

	subs = strings.Split(subs[1], "]")
	if len(subs) != 2 {
		return -1, INVALID, errors.New("malformed array parameter")
	}

	length, err := strconv.Atoi(subs[0])
	if err != nil {
		return -1, INVALID, errors.New("malformed array length")

	}

	gotype := stringToGoType[strings.ToUpper(subs[1])]
	if gotype == INVALID {
		return -1, INVALID, errors.New("malformed array type")
	}

	return length, gotype, nil
}
