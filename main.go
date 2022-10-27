package main

import (
	"fmt"
	"go-lexerkit/lexerkit"
)

func main() {
	type Parser = lexerkit.Parser
	Str := lexerkit.Str
	Seq := lexerkit.Seq
	Alt := lexerkit.Alt
	Reg := lexerkit.Regstr
	Many := lexerkit.Many
	SepBy := lexerkit.SepBy
	TakeWhile := lexerkit.TakeWhile

	lbrace := Str("{")
	rbrace := Str("}")
	lbracket := Str("[")
	rbracket := Str("]")
	comma := Str(",")
	colon := Str(":")
	dquote := Str("\"")

	null := Str("null").Name("Null")
	_true := Str("true").Name("True")
	_false := Str("false").Name("False")
	_string := TakeWhile(func(char rune, _ int) bool {
		return string(char) != "\""
	}).Wrap(dquote, dquote).Name("String")
	number := Reg(`-?(0|[1-9][0-9]*)([.][0-9]+)?([eE][+-]?[0-9]+)?`).Name("Number")
	white := Str(" ")
	optWhite := Many(white)

	value := lexerkit.DummyParser{}

	array := SepBy(value.Parser, comma).Wrap(lbracket, rbracket).Name("Array")
	object := SepBy(Seq(_string.Skip(colon), value.Parser).Wrap(optWhite, optWhite), comma).Wrap(lbrace, rbrace).Name("Object")

	value.InternalParser = Alt(array, object, null, _true, _false, _string, number).Wrap(optWhite, optWhite)

	jsonParse := value.Parser

	target1 := "[{ \"key1\": { \"key2\": [3, 4, 5, [true, false], null] } }]"
	result1, err1 := jsonParse(&target1, 0)
	if err1 != nil {
		fmt.Println(err1)
	} else {
		fmt.Println("success")
	}
	fmt.Print(lexerkit.Stringify(result1))
}
