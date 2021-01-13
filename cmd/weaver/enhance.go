package main

import "fmt"

func checkIfAllTypesAreContainedInGOTIR(g *Gotir) {

	allGood := true
	for _, f := range g.Functions {
		for _, param := range f.Params {

			if !check(param.TypeName, g) {
				allGood = false
				fmt.Printf("Did not find: %s\n", param.TypeName)
			}
		}
	}

	if !allGood {
		fmt.Println("FACK")
	}
}

func check(typeName string, g *Gotir) bool {
	found := false
	for _, b := range g.BaseTypes {
		if typeName == b.Name {
			found = true
			// fmt.Printf("Found %s %d\n", b.Name, b.TypeSize)
		}
	}

	for _, b := range g.PtrTypes {
		if typeName == b.Name {
			found = true
			// fmt.Printf("Found %s\n", b.Name)
		}
	}

	return found
}
