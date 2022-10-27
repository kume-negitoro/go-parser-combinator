package lexerkit

import (
	"fmt"
	"strings"
	"testing"
)

func TestStr(t *testing.T) {
	t.Run("Strパーサがマッチする", func(t *testing.T) {
		target := "test"
		result, err := Str("test")(&target, 0)

		t.Run("エラーでない", func(t *testing.T) {
			if err != nil {
				fmt.Print(Stringify(result))
				t.Fail()
			}
		})

		t.Run(("valueが正しい"), func(t *testing.T) {
			if result.value != target {
				fmt.Print(Stringify(result))
				t.Fail()
			}
		})

		t.Run("indexが正しい", func(t *testing.T) {
			if result.index != len(target) {
				fmt.Print(Stringify(result))
				t.Fail()
			}
		})
	})

	t.Run("Strパーサがマッチしない", func(t *testing.T) {
		target := "test"
		result, err := Str("not test")(&target, 0)

		t.Run("エラーである", func(t *testing.T) {
			if err == nil {
				fmt.Print(Stringify(result))
				t.Fail()
			}
		})

		t.Run("indexが正しい", func(t *testing.T) {
			if result.index != 0 {
				fmt.Print(Stringify(result))
				t.Fail()
			}
		})
	})
}

func TestMany(t *testing.T) {
	t.Run("0回マッチする", func(t *testing.T) {
		target := ""
		result, _ := Many(Str("test"))(&target, 0)

		t.Run("子要素が0個である", func(t *testing.T) {
			if len(result.children) != 0 {
				fmt.Print(Stringify(result))
				t.Fail()
			}
		})

		t.Run("indexが正しい", func(t *testing.T) {
			if result.index != 0 {
				fmt.Print(Stringify(result))
				t.Fail()
			}
		})
	})

	t.Run("2回マッチする", func(t *testing.T) {
		target := strings.Repeat("test", 2)
		result, _ := Many(Str("test"))(&target, 0)

		t.Run("子要素が2個である", func(t *testing.T) {
			if len(result.children) != 2 {
				fmt.Print(Stringify(result))
				t.Fail()
			}
		})

		t.Run("indexが正しい", func(t *testing.T) {
			if result.index != len(target) {
				fmt.Print(Stringify(result))
				t.Fail()
			}
		})
	})
}

func TestSeqMap(t *testing.T) {
	t.Run("", func(t *testing.T) {})
}

func TestSepBy(t *testing.T) {
	t.Run("3回マッチする", func(t *testing.T) {
		target := "a,a,a"
		result, _ := SepBy1(Str("a"), Str(","))(&target, 0)

		t.Run("子要素が3個である", func(t *testing.T) {
			if len(result.children) != 3 {
				fmt.Print(Stringify(result))
				t.Fail()
			}
		})

		t.Run("indexが正しい", func(t *testing.T) {
			if result.index != len(target) {
				fmt.Print(Stringify(result))
				t.Fail()
			}
		})
	})
}

func TestTakeWhile(t *testing.T) {
	t.Run("0文字マッチする", func(t *testing.T) {
		target := "xyz"
		result, _ := TakeWhile(func(char rune, index int) bool {
			return string(char) == "a"
		})(&target, 0)

		t.Run("valueが正しい", func(t *testing.T) {
			if result.value != "" {
				fmt.Print(Stringify(result))
				t.Fail()
			}
		})

		t.Run("indexが正しい", func(t *testing.T) {
			if result.index != 0 {
				fmt.Print(Stringify(result))
				t.Fail()
			}
		})

	})

	t.Run("1文字マッチする", func(t *testing.T) {
		target := "xyz"
		result, _ := TakeWhile(func(char rune, index int) bool {
			return string(char) == "x"
		})(&target, 0)

		t.Run("valueが正しい", func(t *testing.T) {
			if result.value != "x" {
				fmt.Print(Stringify(result))
				t.Fail()
			}
		})

		t.Run("indexが正しい", func(t *testing.T) {
			if result.index != 1 {
				fmt.Print(Stringify(result))
				t.Fail()
			}
		})
	})

	t.Run("終端までマッチする", func(t *testing.T) {
		target := "aaa"
		result, _ := TakeWhile(func(char rune, index int) bool {
			return string(char) == "a"
		})(&target, 0)

		t.Run("valueが正しい", func(t *testing.T) {
			if result.value != "aaa" {
				fmt.Print(Stringify(result))
				t.Fail()
			}
		})

		t.Run("indexが正しい", func(t *testing.T) {
			if result.index != len(target) {
				fmt.Print(Stringify(result))
				t.Fail()
			}
		})
	})
}

func TestSkip(t *testing.T) {
	t.Run("要素をスキップできる", func(t *testing.T) {
		target := "a,"
		result, err := Str("a").Skip(Str(","))(&target, 0)

		t.Run("エラーでない", func(t *testing.T) {
			if err != nil {
				fmt.Print(Stringify(result))
				t.Fail()
			}
		})

		t.Run("valueが正しい", func(t *testing.T) {
			if result.value != "a" {
				fmt.Print(Stringify(result))
				t.Fail()
			}
		})

		t.Run("indexが正しい", func(t *testing.T) {
			if result.index != len(target) {
				fmt.Print(Stringify(result))
				t.Fail()
			}
		})
	})
}

func TestRegExp(t *testing.T) {
	t.Run("1つの空白にマッチする", func(t *testing.T) {
		target := " "
		result, err := Regstr(`%s`)(&target, 0)

		t.Run("エラーでない", func(t *testing.T) {
			if err != nil {
				fmt.Println(err)
				fmt.Print(Stringify(result))
				t.Fail()
			}
		})

		t.Run("valueが正しい", func(t *testing.T) {
			if result.value != target {
				fmt.Println(err)
				fmt.Print(Stringify(result))
				t.Fail()
			}
		})

		t.Run("indexが正しい", func(t *testing.T) {
			if result.index != len(target) {
				fmt.Println(err)
				fmt.Print(Stringify(result))
				t.Fail()
			}
		})
	})
}
