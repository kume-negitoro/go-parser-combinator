package lexerkit

import (
	"fmt"
	"regexp"
)

// パーサ型
type Parser func(target *string, index int) (*Result, error)

// パーサを加工する関数の型
type ParserMapper func(parser Parser) Parser

// パーサを遅延評価するための構造
type DummyParser struct {
	InternalParser Parser
}

// 内部のパーサを使うパーサ
func (dummy *DummyParser) Parser(target *string, index int) (*Result, error) {
	return dummy.InternalParser(target, index)
}

// 必ず成功し、与えたvalueを持ったResultを返すパーサ
func Succeed(value string) Parser {
	return func(target *string, index int) (*Result, error) {
		return MakeSuccess(index, value), nil
	}
}

// 必ず失敗し、与えたvalueを持ったResultを返すパーサ
func Failed(expected []string) Parser {
	return func(target *string, index int) (*Result, error) {
		return MakeFailure(index, expected),
			fmt.Errorf("ParseError: %s is expected", expected)
	}
}

// 与えられた文字列にマッチするパーサ
func Str(expected string) Parser {
	return func(target *string, index int) (*Result, error) {
		length := len(expected)

		if index+length > len(*target) {
			return MakeFailure(index, []string{expected}),
				fmt.Errorf("ParseError: %s is expected at index==%d", expected, index)
		}

		if (*target)[index:index+length] == expected {
			return MakeSuccess(index+length, expected), nil
		} else {
			return MakeFailure(index, []string{expected}),
				fmt.Errorf("ParseError: %s is expected at index==%d", expected, index)
		}
	}
}

// 複数回マッチするパーサ
func Many(parser Parser) Parser {
	return func(target *string, index int) (*Result, error) {
		children := []*Result{}
		for {
			parsed, err := parser(target, index)
			if err != nil {
				break
			}

			children = append(children, parsed)
			index = parsed.index
		}

		return MakeContainer(index, children), nil
	}
}

// いずれかのパーサがマッチするパーサ
func Alt(parsers ...Parser) Parser {
	return func(target *string, index int) (*Result, error) {
		expected := make([]string, 0, len(parsers))

		for _, parser := range parsers {
			parsed, err := parser(target, index)
			if err == nil {
				return parsed, nil
			}

			expected = append(expected, fmt.Sprintf("%s", parsed.expected))
		}

		return MakeFailure(index, expected),
			fmt.Errorf("ParseError: One of %s is expected", expected)
	}
}

// どちらかのパーサがマッチするパーサ
func Or(left Parser, right Parser) Parser {
	return Alt(left, right)
}

// マッチしてもしなくても良いパーサ
func Opt(parser Parser) Parser {
	return Or(parser, Succeed(""))
}

// 与えられた順のパーサが全てマッチするパーサ
func Seq(parsers ...Parser) Parser {
	return func(target *string, index int) (*Result, error) {
		children := make([]*Result, 0, len(parsers))

		for _, parser := range parsers {
			parsed, err := parser(target, index)
			if err != nil {
				return MakeFailure(index, parsed.expected), err
			}

			children = append(children, parsed)
			index = parsed.index
		}

		return MakeContainer(index, children), nil
	}
}

// 与えられた正規表現にマッチするパーサ
func Regexp(reg *regexp.Regexp) Parser {
	return func(target *string, index int) (*Result, error) {
		matched := reg.FindString((*target)[index:])

		if matched != "" {
			return MakeSuccess(index+len(matched), matched), nil
		} else {
			return MakeFailure(index, []string{reg.String()}),
				fmt.Errorf("ParseError: A string matching /%s/ is expected", reg.String())
		}
	}
}

// 与えられた正規表現文字列にマッチするパーサ
func Regstr(regstr string) Parser {
	return Regexp(regexp.MustCompile(regstr))
}

// Resultをmapperによって変更するパーサ
func Map(mapper ResultMapper, parser Parser) Parser {
	return func(target *string, index int) (*Result, error) {
		return mapper(parser(target, index))
	}
}

// SeqのResultをmapperによって変更するパーサ
func SeqMap(mapper ResultMapper, parsers ...Parser) Parser {
	return Map(mapper, Seq(parsers...))
}

// セパレータで区切られた部分を少なくとも1回パースするパーサ
func SepBy1(parser Parser, separator Parser) Parser {
	pairs := Many(Seq(separator, parser))
	return SeqMap(func(result *Result, err error) (*Result, error) {
		if err != nil {
			return result, err
		}

		resultNum := len(result.children)
		results := make([]*Result, 0, resultNum)
		results = append(results, result.children[0])

		for _, result := range result.children[1].children {
			results = append(results, result.children[1])
		}

		return MakeContainer(results[len(results)-1].index, results), nil
	}, parser, pairs)
}

// セパレータで区切られた部分をパースするパーサ
func SepBy(parser Parser, separator Parser) Parser {
	return Opt(SepBy1(parser, separator))
}

// テスト関数が通る間パースするパーサ
func TakeWhile(test func(char rune, index int) bool) Parser {
	return func(target *string, index int) (*Result, error) {
		for i, c := range (*target)[index:] {
			if !test(c, index) {
				return MakeSuccess(index+i, (*target)[index:index+i]), nil
			}
		}
		return MakeSuccess(len(*target), (*target)[index:]), nil
	}
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
	return Opt(parser)
}

// Resultを書き換える
func (parser Parser) Map(mapper ResultMapper) Parser {
	return Map(mapper, parser)
}

// パーサを加工する
func (parser Parser) Thru(wrapper ParserMapper) Parser {
	return wrapper(parser)
}

// leftとrightに囲まれた部分のResultを返す
func (parser Parser) Wrap(left Parser, right Parser) Parser {
	return SeqMap(func(result *Result, err error) (*Result, error) {
		if err != nil {
			return result, err
		}

		return MergeResults(result.children[1], result.children[2]), nil
	}, left, parser, right)
}

// 次のResultを返す
func (parser Parser) Then(next Parser) Parser {
	return SeqMap(func(result *Result, err error) (*Result, error) {
		if err != nil {
			return result, err
		}

		return result.children[1], nil
	}, parser, next)
}

// 次のResultを飛ばす
func (parser Parser) Skip(next Parser) Parser {
	return SeqMap(func(result *Result, err error) (*Result, error) {
		if err != nil {
			return result, err
		}

		return MergeResults(result.children[0], result.children[1]), nil
	}, parser, next)
}

// 名前をつける
func (parser Parser) Name(name string) Parser {
	return Map(func(result *Result, err error) (*Result, error) {
		if err != nil {
			return result, err
		}

		result.name = name
		return result, nil
	}, parser)
}
