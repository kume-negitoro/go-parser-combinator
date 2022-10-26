package main

import (
	"fmt"
	"go-lexerkit/lexerkit"
)

func main() {
	Str := lexerkit.Str
	Seq := lexerkit.Seq
	Alt := lexerkit.Alt
	// Reg := lexerkit.Regstr

	lbrace := Str("{")
	rbrace := Str("}")
	lbracket := Str("(")
	rbracket := Str(")")
	// comma := Str(",")
	// colon := Str(":")
	null := Str("null")
	_true := Str("true")
	_false := Str("false")
	// key := Reg(`".+?":`)
	key := Str("test:")
	value := Alt(null, _true, _false)

	jsonParse := Alt(value, Seq(lbrace, key, value, rbrace), Seq(lbracket, value, rbracket))

	target := "{test:null}"
	result1 := jsonParse(&target, 0)
	fmt.Print(lexerkit.Stringify(result1))

}
