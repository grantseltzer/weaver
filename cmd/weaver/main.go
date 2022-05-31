package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {

	bc, err := init_binary_context(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(bc.Functions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}
