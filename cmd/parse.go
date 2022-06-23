package main

import (
	"debug/dwarf"
	"debug/elf"
	"errors"
	"fmt"
)

type binary_context struct {
	elf_file          *elf.File
	text_section_addr uint64
	syms_to_offset    map[string]uint64
	Functions         []*function_info
}

type function_info struct {
	Name       string
	TextOffset uint64
	Params     []function_parameter
}

type function_parameter struct {
	Name     string
	TypeName string
	TypeSize int64 // how much space it takes on the stack, in bytes
	IsReturn bool

	//TODO:
	ParamNumber uint32
	IsOnStack   bool
	StackOffset uint32
}

//FIXME: this is messed up, need to fix (replace helpers function in main)
func init_binary_context(path string) (*binary_context, error) {
	e := &binary_context{
		syms_to_offset: make(map[string]uint64),
	}
	f, err := elf.Open(path)
	if err != nil {
		return nil, err
	}
	e.elf_file = f

	// Collect symbols offset within ELF .text section
	textSection := e.elf_file.Section(".text")
	if textSection == nil {
		return nil, errors.New("no .text section")
	}
	e.text_section_addr = textSection.Addr + textSection.Offset

	syms, err := e.elf_file.Symbols()
	if err != nil {
		return nil, err
	}

	for _, sym := range syms {
		e.syms_to_offset[sym.Name] = sym.Value - e.text_section_addr
	}

	// Collect function information from DWARF
	err = e.parse_dwarf_function_info()
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (e *binary_context) parse_dwarf_function_info() error {
	data, err := e.elf_file.DWARF()
	if err != nil {
		return err
	}

	lineReader := data.Reader()
	typeReader := data.Reader()
	dwarfEntryIndex := map[string]*dwarf.Entry{}

	func_info := []*function_info{}

	currentlyReadingFunction := &function_info{}

entryReadLoop:
	for {
		entry, err := lineReader.Next()
		if err == nil && entry == nil {
			break
		}
		if err != nil {
			return err
		}

		if entryIsNull(entry) {
			if currentlyReadingFunction != nil {
				func_info = append(func_info, currentlyReadingFunction)
				currentlyReadingFunction = nil
			}
			continue entryReadLoop
		}

		// Index all possible types by name for later lookup
		for _, field := range entry.Field {
			if field.Attr == dwarf.AttrName {
				dwarfEntryIndex[field.Val.(string)] = entry
			}
		}

		// Found a function
		if entry.Tag == dwarf.TagSubprogram {
			currentlyReadingFunction = e.read_function_init(entry)
		}
		// If currently reading the parameters of a function
		if currentlyReadingFunction != nil && entry.Tag == dwarf.TagFormalParameter {
			err = read_function_parameter(typeReader, entry, currentlyReadingFunction)
			if err != nil {
				return err
			}
		}
	}

	e.Functions = func_info

	return nil
}

func (e *binary_context) symbol_name_to_offset(symbol string) (uint64, error) {
	offset, ok := e.syms_to_offset[symbol]
	if !ok {
		return 0, fmt.Errorf("could not find symbol: %s", symbol)
	}
	return offset, nil
}

func (e *binary_context) read_function_init(entry *dwarf.Entry) *function_info {
	currentlyReadingFunction := &function_info{}

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

	offset, err := e.symbol_name_to_offset(currentlyReadingFunction.Name)
	if err != nil {
		return nil
	}
	currentlyReadingFunction.TextOffset = offset
	currentlyReadingFunction.Params = []function_parameter{}

	return currentlyReadingFunction
}

func read_function_parameter(typeReader *dwarf.Reader, entry *dwarf.Entry, currentlyReadingFunction *function_info) error {

	var (
		typeEntry *dwarf.Entry
		err       error
	)

	newParam := function_parameter{IsReturn: false}
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
