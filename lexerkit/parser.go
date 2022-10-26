package lexerkit

import (
	"fmt"
	"regexp"
)

type Parser func(target string, index int) *Result

// パーサを遅延評価で読み込むパーサ
func Lazy(fn func() Parser) Parser {
	return func(target string, index int) *Result {
		return fn()(target, index)
	}
}

// 必ず成功し、与えたvalueを持ったResultを返すパーサ
func Succeed(value string) Parser {
	return func(target string, index int) *Result {
		return MakeSuccess(index, value)
	}
}

// 必ず失敗し、与えたvalueを持ったResultを返すパーサ
func Failed(expected []string) Parser {
	return func(target string, index int) *Result {
		return MakeFailure(index, expected)
	}
}

// 与えられた文字列にマッチするパーサ
func Str(expected string) Parser {
	return func(target string, index int) *Result {
		var length = len(expected)

		if index+length > len(target) {
			return MakeEmpty(false, index)
		}

		if target[index:index+length] == expected {
			return MakeSuccess(index+length, expected)
		} else {
			return MakeFailure(index, []string{"'" + expected + "'"})
		}
	}
}

// 複数回マッチするパーサ
func Many(parser Parser) Parser {
	return func(target string, index int) *Result {
		var children = []*Result{}
		for {
			var parsed = parser(target, index)

			if parsed.Status() {
				children = append(children, parsed)
				index = parsed.Index()
			} else {
				break
			}
		}

		return MakeContainer(index, children)
	}
}

// いずれかのパーサがマッチするパーサ
func Alt(parsers ...Parser) Parser {
	return func(target string, index int) *Result {
		for _, parser := range parsers {
			var parsed = parser(target, index)

			if parsed.Status() {
				fmt.Printf("debug: ok, length: %d\n", len(parsers))
				return parsed
			}
			fmt.Println("debug: ng")
		}

		return MakeFailure(index, []string{"no any parsers matched"})
	}
}

// どちらかのパーサがマッチするパーサ
func Or(left Parser, right Parser) Parser {
	return Alt(left, right)
}

// 与えられた順のパーサが全てマッチするパーサ
func Seq(parsers ...Parser) Parser {
	return func(target string, index int) *Result {
		var children = []*Result{}

		for _, parser := range parsers {
			var parsed = parser(target, index)

			if parsed.Status() {
				children = append(children, parsed)
				index = parsed.Index()
			} else {
				return MakeFailure(index, parsed.Expected())
			}
		}

		return MakeContainer(index, children)
	}
}

// 与えられた正規表現にマッチするパーサ
func Regexp(reg *regexp.Regexp) Parser {
	return func(target string, index int) *Result {
		var matched = reg.FindString(target[index:])

		if matched != "" {
			return MakeSuccess(index+len(matched), matched)
		} else {
			return MakeFailure(index, []string{"regexp missmatch"})
		}
	}
}

// 与えられた正規表現文字列にマッチするパーサ
func Regstr(regstr string) Parser {
	return Regexp(regexp.MustCompile(regstr))
}

// 複数回マッチする
func (parser Parser) Many() Parser {
	return Many(parser)
}

// いずれかにマッチする
func (parser Parser) Alt(parsers ...Parser) Parser {
	return Alt(append([]Parser{parser}, parsers...)...)
}

// どちらかにマッチする
func (parser Parser) Or(alt Parser) Parser {
	return Or(parser, alt)
}

// 与えられた順にマッチする
func (parser Parser) Seq(parsers ...Parser) Parser {
	return Seq(append([]Parser{parser}, parsers...)...)
}

// マッチしなくても良い
func (parser Parser) Opt() Parser {
	return Or(parser, Succeed(""))
}
