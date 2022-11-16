# go-parser-combinator
これはGoの勉強として実装したシンプルなパーサコンビネータです。

## パーサ

パーサは`*string`型の`target`と`int`型の`index`を受け取り、`*Result`型と`error`型を返す関数です。

`target`の`index`から始まる部分文字列に対してパーサが正常に解析できた場合は、`status`が`true`である`*Result`が返却され、`error`は`nil`となります。正常に解析ができなかった場合は、`status`が`false`である`*Result`が返却され、`error`は`ParseError`となります。

パーサは以下のtypeで表されます。

```go
type Parser func(target *string, index int) (*Result, error)
```

### Succeed
必ず成功し、与えたvalueを持ったResultを返すパーサ
```go
func Succeed(value string) Parser
```

### Failed
必ず失敗し、与えたvalueを持ったResultを返すパーサ
```go
func Failed(expected []string) Parser
```

### Str
与えられた文字列にマッチするパーサ
```go
func Str(expected string) Parser
```

### Many
与えられたパーサが複数回マッチするパーサ
```go
func Many(parser Parser) Parser
```

### Alt
与えられたいずれかのパーサがマッチするパーサ
```go
func Alt(parsers ...Parser) Parser
```

### Or
与えられた2つのパーサのうち、どちらかのパーサがマッチするパーサ
```go
func Or(left Parser, right Parser) Parser
```

### Opt
マッチの結果に関わらず成功するパーサ
```go
func Opt(parser Parser) Parser
```

### Seq
与えられた順のパーサが全てマッチするパーサ
```go
func Seq(parsers ...Parser) Parser
```

### Regexp
与えられた正規表現にマッチするパーサ
```go
func Regexp(reg *regexp.Regexp) Parser
```

### Regstr
与えられた正規表現文字列にマッチするパーサ
```go
func Regstr(regstr string) Parser
```

### Map
Resultをmapperによって変更するパーサ
```go
func Map(mapper ResultMapper, parser Parser) Parser
```

### SeqMap
SeqのResultをmapperによって変更するパーサ
```go
func SeqMap(mapper ResultMapper, parsers ...Parser) Parser
```

### SepBy1
セパレータで区切られた部分を少なくとも1回パースするパーサ
```go
func SepBy1(parser Parser, separator Parser) Parser
```

### SepBy
セパレータで区切られた部分をパースするパーサ
```go
func SepBy(parser Parser, separator Parser) Parser
```

### TakeWhile
テスト関数が真の間パースするパーサ
```go
func TakeWhile(test func(char rune, index int) bool) Parser
```

### Thru (メソッド)
パーサを加工するパーサを返す
```go
func (parser Parser) Thru(wrapper ParserMapper) Parser {
	return wrapper(parser)
}
```

### Wrap (メソッド)
leftとrightに囲まれた部分のResultを返すパーサを返す
```go
func (parser Parser) Wrap(left Parser, right Parser) Parser
```

### Then (メソッド)
nextの部分のResultを返すパーサを返す
```go
func (parser Parser) Then(next Parser) Parser
```

### Skip (メソッド)
nextの部分のResultを無視するパーサを返す
```go
func (parser Parser) Skip(next Parser) Parser
```

### Name (メソッド)
Resultに名前を付与するパーサを返す
```go
func (parser Parser) Name(name string) Parser
```
