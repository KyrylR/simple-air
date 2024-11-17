package main

import (
	"log"
	"strconv"

	"github.com/consensys/gnark-crypto/field/generator"
	field "github.com/consensys/gnark-crypto/field/generator/config"
)

func main() {
	elementName := "Element"

	fIntegration, err := field.NewFieldConfig("ff", elementName, strconv.FormatUint(Modulus, 10), false)
	if err != nil {
		log.Fatal(elementName, err)
	}

	if err = generator.GenerateFF(fIntegration, "./"); err != nil {
		log.Fatal(elementName, err)
	}
}
