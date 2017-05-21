package main

import (
	"fmt"

	"github.com/damnever/cc"
)

func pp(c cc.Configer) {
	fmt.Println("---")
	fmt.Println(c.String("name"))
	cc := c.Config("map")
	fmt.Println(cc.Bool("key_one"))

	child := cc.Value("child").Map()
	fmt.Println(child["key_three"].Int())
	fmt.Println(child["key_four"].String())

	list := c.Value("list").List()
	fmt.Println(list[0].String())
	fmt.Println(list[1].Int())
	fmt.Println(list[2].Float())
	fmt.Println(list[3].Bool())

	patterns := c.Config("patterns")
	fmt.Println(patterns.Pattern("string_pattern").ValidateString("aaaaa"))
	fmt.Println(patterns.Pattern("int_pattern").ValidateInt(3))
	fmt.Println(patterns.Pattern("float_pattern").ValidateFloat(0.5))
}

func main() {
	cj, err := cc.NewConfigFromFile("example.json")
	must(err)
	pp(cj)

	cy, err := cc.NewConfigFromFile("example.yaml")
	must(err)
	pp(cy)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
