package main

import (
	"debug/dwarf"
	"debug/elf"
	"fmt"
	"log"
	"os"
)

func maain() {
	elfFile, err := elf.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	dwarfData, err := elfFile.DWARF()
	if err != nil {
		log.Fatal(err)
	}

	dwarfReader := dwarfData.Reader()
	// otherReader := dwarfData.Reader()

	nextIsParam := false

	for {
		entry, err := dwarfReader.Next()
		if err != nil {
			log.Fatal(err)
		}

		if nextIsParam {
			fmt.Println(entry)
			fmt.Println(entryIsNull(entry))
			os.Exit(0)
		}

		if entry.Tag == dwarf.TagSubprogram {
			for i := range entry.Field {
				if entry.Field[i].Attr == dwarf.AttrName {
					fmt.Println(entry.Field[i].Val.(string))
				}
			}
			nextIsParam = true
		}

		if entry == nil {
			os.Exit(2)
		}
	}
}

func entryIsNull(e *dwarf.Entry) bool {
	return e.Children == false &&
		len(e.Field) == 0 &&
		e.Offset == 0 &&
		e.Tag == dwarf.Tag(0)
}
