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

	null := Str("null")
	_true := Str("true")
	_false := Str("false")
	_string := TakeWhile(func(char rune, _ int) bool {
		return string(char) != "\""
	}).Wrap(dquote, dquote)
	number := Reg(`-?(0|[1-9][0-9]*)([.][0-9]+)?([eE][+-]?[0-9]+)?`)
	white := Str(" ")
	optWhite := Many(white)

	value := lexerkit.DummyParser{}

	array := SepBy(value.Parser, comma).Wrap(lbracket, rbracket)
	object := SepBy(Seq(_string.Skip(colon), value.Parser).Wrap(optWhite, optWhite), comma).Wrap(lbrace, rbrace)

	value.InternalParser = Alt(array, object, null, _true, _false, _string, number).Wrap(optWhite, optWhite)

	jsonParse := value.Parser

	target1 := "[{ \"key1\": { \"key2\": [] } }]"
	result1, err1 := jsonParse(&target1, 0)
	if err1 != nil {
		fmt.Println(err1)
	} else {
		fmt.Println("success")
	}
	fmt.Print(lexerkit.Stringify(result1))
}
