package main

import (
	"debug/dwarf"
	"debug/elf"
	"strings"
)

type Gotir struct {
	Functions []*function_type
}

type function_type struct {
	Name   string
	Params []function_param
}

type function_param struct {
	Name           string
	TypeName       string
	StartingOffset uint
	TypeSize       int64 // how much space it takes on the stack, in bytes
	IsReturn       bool
}

var entryIndex = map[string]*dwarf.Entry{}

func getSizesAndStackOffsets(ir *Gotir, data *dwarf.Data) {
	// TODO:
	// Go through each function in ir.Functions
	// Iterate over the params, having a dwarf reader go and recursive determine the size on the stack of each
	// parameter and their starting offsets

	for _, f := range ir.Functions {

		for _, param := range f.Params {
			// Use TypeName to find TypeSize, then afterwards calculate StartingOffset
			entry := entryIndex[param.TypeName]
			if entry == nil {
				continue
			}

			// Look for size
			for i := range entry.Field {
				if entry.Field[i].Attr == dwarf.AttrByteSize {
					param.TypeSize = entry.Field[i].Val.(int64)
				}
			}

			if param.TypeSize == 0 && strings.HasPrefix(param.Name, "*") {
				param.TypeSize = 8
			}

			// TODO: Otherwise just don't support it for now

		}
	}
}

// parseFromPath reads in all of the type information from the DWARF section of the ELF at the given patho
func parseFromPath(path string) (*Gotir, error) {
	elfFile, err := elf.Open(path)
	if err != nil {
		return nil, err
	}

	data, err := elfFile.DWARF()
	if err != nil {
		return nil, err
	}

	ir, err := parseFromData(data)
	if err != nil {
		return nil, err
	}

	getSizesAndStackOffsets(ir, data)
	return ir, nil
}

func parseFromData(data *dwarf.Data) (*Gotir, error) {

	lineReader := data.Reader()
	typeReader := data.Reader()

	ir := &Gotir{
		Functions: []*function_type{},
	}

	var currentlyReadingFunction *function_type = nil

entryReadLoop:
	for {
		entry, err := lineReader.Next()
		if err == nil && entry == nil {
			break
		}
		if err != nil {
			return nil, err
		}

		if entryIsNull(entry) {
			if currentlyReadingFunction != nil {
				ir.Functions = append(ir.Functions, currentlyReadingFunction)
				currentlyReadingFunction = nil
			}
			continue entryReadLoop
		}

		// Index all possible types by name for later lookup
		for _, field := range entry.Field {
			if field.Attr == dwarf.AttrName {
				entryIndex[field.Val.(string)] = entry
			}
		}

		// Found a function
		if entry.Tag == dwarf.TagSubprogram {
			currentlyReadingFunction = readFunctionInit(entry)
		}
		// If currently reading the parameters of a function
		if currentlyReadingFunction != nil && entry.Tag == dwarf.TagFormalParameter {
			err = readFunctionParameter(typeReader, entry, currentlyReadingFunction)
			if err != nil {
				return nil, err
			}
		}
	}

	return ir, nil
}

func readFunctionInit(entry *dwarf.Entry) *function_type {
	currentlyReadingFunction := &function_type{}

	isNamedSubroutine := false
	for _, field := range entry.Field {
		if field.Attr == dwarf.AttrName {
			isNamedSubroutine = true
			currentlyReadingFunction.Name = field.Val.(string)
		}
	}

	if !isNamedSubroutine {
		return nil
	}

	currentlyReadingFunction.Params = []function_param{}
	return currentlyReadingFunction
}

func readFunctionParameter(typeReader *dwarf.Reader, entry *dwarf.Entry, currentlyReadingFunction *function_type) error {

	var (
		typeEntry *dwarf.Entry
		err       error
	)

	newParam := function_param{IsReturn: false}
	isNamedParameter := false
	for _, field := range entry.Field {

		if field.Attr == dwarf.AttrName {
			newParam.Name = field.Val.(string)
			isNamedParameter = true
		}

		if field.Attr == dwarf.AttrVarParam {
			newParam.IsReturn = field.Val.(bool)
		}

		// Get the name of the type of the parameter
		// XXX: Have to go back later to get the size
		if field.Attr == dwarf.AttrType {
			typeReader.Seek(field.Val.(dwarf.Offset))
			typeEntry, err = typeReader.Next()
			if err != nil {
				return err
			}
			for i := range typeEntry.Field {
				if typeEntry.Field[i].Attr == dwarf.AttrName {
					newParam.TypeName = typeEntry.Field[i].Val.(string)
				}
			}
		}
	}

	if isNamedParameter {
		currentlyReadingFunction.Params = append(currentlyReadingFunction.Params, newParam)
	}
	return nil
}

func entryIsNull(e *dwarf.Entry) bool {
	return e.Children == false &&
		len(e.Field) == 0 &&
		e.Offset == 0 &&
		e.Tag == dwarf.Tag(0)
}
