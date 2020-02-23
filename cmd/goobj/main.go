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

	/* Line Table Info */
	// lineTableData, err := elfFile.Section(".gopclntab").Data()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// addr := elfFile.Section(".text").Addr
	// lineTable := gosym.NewLineTable(lineTableData, addr)
	// symTable, err := gosym.NewTable([]byte{}, lineTable)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for _, funcer := range symTable.Funcs {
	// 	// fmt.Printf("0x%x\t%s\n", funcer.Value, funcer.BaseName())
	// }

	/* DWARF info */
	dwarfData, err := elfFile.DWARF()
	if err != nil {
		log.Fatal(err)
	}

	reader := dwarfData.Reader()

	for {
		entry, err := reader.Next()
		if err != nil {
			continue
		}

		if entry == nil {
			os.Exit(2)
		}

		t, err := dwarfData.Type(entry.Offset)
		if err != nil {
			continue
		}

		switch x := t.(type) {

		case *dwarf.FuncType:

			fmt.Println(entry.AttrField(dwarf.AttrName))
			fmt.Printf("Offset: %x\tFunc: %s \n", entry.Offset, x.String())
			for _, f := range entry.Field {
				fmt.Println(f)
			}
			fmt.Printf("\n\n")
		}
	}
}

/*
Try:
Get PC range from Data.Ranges

*/
