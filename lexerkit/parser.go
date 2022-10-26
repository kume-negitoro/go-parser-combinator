package lexerkit

import (
	"regexp"
)

type Parser func(target *string, index int) *Result

// パーサを遅延評価で読み込むパーサ
func Lazy(fn func() Parser) Parser {
	return func(target *string, index int) *Result {
		return fn()(target, index)
	}
}

// 必ず成功し、与えたvalueを持ったResultを返すパーサ
func Succeed(value string) Parser {
	return func(target *string, index int) *Result {
		return MakeSuccess(index, value)
	}
}

// 必ず失敗し、与えたvalueを持ったResultを返すパーサ
func Failed(expected []string) Parser {
	return func(target *string, index int) *Result {
		return MakeFailure(index, expected)
	}
}

// 与えられた文字列にマッチするパーサ
func Str(expected string) Parser {
	return func(target *string, index int) *Result {
		var length = len(expected)

		if index+length > len(*target) {
			return MakeEmpty(false, index)
		}

		if (*target)[index:index+length] == expected {
			return MakeSuccess(index+length, expected)
		} else {
			return MakeFailure(index, []string{"'" + expected + "'"})
		}
	}
}

// 複数回マッチするパーサ
func Many(parser Parser) Parser {
	return func(target *string, index int) *Result {
		var children = []*Result{}
		for {
			var parsed = parser(target, index)

			if parsed.status {
				children = append(children, parsed)
				index = parsed.index
			} else {
				break
			}
		}

		return MakeContainer(index, children)
	}
}

// いずれかのパーサがマッチするパーサ
func Alt(parsers ...Parser) Parser {
	return func(target *string, index int) *Result {
		for _, parser := range parsers {
			var parsed = parser(target, index)

			if parsed.status {
				return parsed
			}
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
	return func(target *string, index int) *Result {
		var children = []*Result{}

		for _, parser := range parsers {
			var parsed = parser(target, index)

			if parsed.status {
				children = append(children, parsed)
				index = parsed.index
			} else {
				return MakeFailure(index, parsed.expected)
			}
		}

		return MakeContainer(index, children)
	}
}

// 与えられた正規表現にマッチするパーサ
func Regexp(reg *regexp.Regexp) Parser {
	return func(target *string, index int) *Result {
		var matched = reg.FindString((*target)[index:])

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

// Resultをmapperによって変更するパーサ
func Map(mapper func(result *Result) *Result, parser Parser) Parser {
	return func(target *string, index int) *Result {
		return mapper(parser(target, index))
	}
}

// SeqのResultをmapperによって変更するパーサ
func SeqMap(mapper func(result *Result) *Result, parsers ...Parser) Parser {
	return Map(mapper, Seq(parsers...))
}

// セパレータで区切られた部分をパースするパーサ
func SepBy1(parser Parser, separator Parser) Parser {
	pairs := Many(Seq(separator, parser))
	return SeqMap(func(result *Result) *Result {
		resultNum := len(result.children)
		results := make([]*Result, 0, resultNum)
		for _, result := range result.children {
			results = append(results, result.children[1])
		}
		return MakeContainer(len(result.children), results)
	}, parser, pairs)
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

// Resultを書き換える
func (parser Parser) Map(mapper func(result *Result) *Result) Parser {
	return Map(mapper, parser)
}

// パーサを加工する
func (parser Parser) Thru(wrapper func(parser Parser) Parser) Parser {
	return wrapper(parser)
}

// leftとrightに囲まれた部分のResultを返す
func (parser Parser) Wrap(left Parser, right Parser) Parser {
	return SeqMap(func(result *Result) *Result {
		return result.children[1]
	}, left, parser, right)
}

// 次のResultを返す
func (parser Parser) Then(next Parser) Parser {
	return SeqMap(func(result *Result) *Result {
		return result.children[1]
	}, parser, next)
}

// 次のResultを飛ばす
func (parser Parser) Skip(next Parser) Parser {
	return SeqMap(func(result *Result) *Result {
		return result.children[0]
	}, parser, next)
}
