package main

import (
	"fmt"

	"github.com/fatih/color"
)

func main() {
	color.Cyan("abc")
	fmt.Println(color.RedString("abc"))
}
