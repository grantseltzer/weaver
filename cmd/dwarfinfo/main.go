package main

import (
	"debug/dwarf"
	"debug/elf"
	"fmt"
	"log"
	"os"
)

func main() {
	elfFile, err := elf.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	dwarfData, err := elfFile.DWARF()
	if err != nil {
		log.Fatal(err)
	}

	dwarfReader := dwarfData.Reader()

	otherReader := dwarfData.Reader()

	nextIsParam := false

	for {
		entry, err := dwarfReader.Next()
		if err != nil {
			log.Fatal(err)
		}

		if nextIsParam {
			for _, field := range entry.Field {
				if field.Attr == dwarf.AttrType {
					otherReader.Seek(field.Val.(dwarf.Offset))
					entry, err := otherReader.Next()
					if err != nil {
						log.Fatal("wtf", err)
					}
					fmt.Println(entry)
				}
			}

			nextIsParam = false
		}

		if entry.Tag == dwarf.TagSubprogram {
			nextIsParam = true
		}

		if entry == nil {
			os.Exit(2)
		}
	}
}
