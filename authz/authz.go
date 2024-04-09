package main

import (
	_ "embed"
	"fmt"
)

//go:embed schema.perm
var schema string

func main() {
	// Initialize the schema.
	fmt.Println(schema)
}
