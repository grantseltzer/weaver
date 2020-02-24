package main

import (
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

	syms, _ := elfFile.Symbols()

	for _, sym := range syms {
		if sym.Name == "main.verifyEntitlementArgs" {
			fmt.Printf("%+v\n\n\n\n", sym)
		}
	}

	dwarfData, err := elfFile.DWARF()
	if err != nil {
		log.Fatal(err)
	}

	dwarfReader := dwarfData.Reader()

	for {

		entry, err := dwarfReader.Next()
		if err != nil {
			log.Fatal(err)
		}

		if entry == nil {
			os.Exit(2)
		}

		fmt.Printf("%+v\n\n", entry)

	}

}
