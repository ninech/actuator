package main

import (
	"fmt"

	"github.com/ninech/frau-schultz/api"
)

func main() {
	fmt.Println("Hallo, ich bin Frau Schultz!")

	engine := api.GetMainEngine()
	engine.Run(":8080")
}
