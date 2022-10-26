package lexerkit

import (
	"strconv"
	"strings"
)

type ResultTypes string

const (
	Container ResultTypes = "container"
	Content   ResultTypes = "content"
)

type Result struct {
	resultType ResultTypes
	status     bool
	index      int
	expected   []string
	children   []*Result
	value      string
}

func NewResult() *Result {
	return &Result{
		resultType: Content,
		status:     true,
		index:      0,
		expected:   []string{},
		children:   []*Result{},
		value:      "",
	}
}

func (p *Result) IsEmpty() bool {
	switch p.resultType {
	case Container:
		return len(p.children) == 0
	case Content:
		return p.value == ""
	default:
		return true
	}
}

func MakeContent(
	index int,
	expected []string,
	children []*Result,
	value string,
) *Result {
	return &Result{
		resultType: "content",
		status:     true,
		index:      index,
		expected:   expected,
		children:   []*Result{},
		value:      value,
	}
}

func MakeContainer(
	index int,
	children []*Result,
) *Result {
	return &Result{
		resultType: "container",
		status:     true,
		index:      index,
		expected:   []string{},
		children:   children,
		value:      "",
	}
}

func MakeEmpty(status bool, index int) *Result {
	return &Result{
		resultType: "content",
		status:     status,
		index:      index,
		expected:   []string{},
		children:   []*Result{},
		value:      "",
	}
}

func MakeSuccess(index int, value string) *Result {
	return &Result{
		resultType: "content",
		status:     true,
		index:      index,
		expected:   []string{},
		children:   []*Result{},
		value:      value,
	}
}

func MakeFailure(index int, expected []string) *Result {
	return &Result{
		resultType: "content",
		status:     false,
		index:      index,
		expected:   expected,
		children:   []*Result{},
		value:      "",
	}
}

func MergeResults(result Result, last Result) *Result {
	if result.index > last.index {
		return &result
	}

	return &Result{
		resultType: result.resultType,
		status:     result.status,
		index:      last.index,
		expected:   result.expected,
		children:   result.children,
		value:      result.value,
	}
}

type Success struct{ Result }
type Failure struct{ Result }

func Stringify(result *Result) string {
	var tab = func(buffer string, n int) string {
		return buffer + strings.Repeat("  ", n)
	}
	var loop func(result *Result, buffer string, nest int) string
	loop = func(result *Result, buffer string, nest int) string {
		buffer = tab(buffer, nest)
		buffer = buffer + "{\n"

		buffer = tab(buffer, nest+1)
		buffer = buffer + "\"status\": " + strconv.FormatBool(result.status) + ",\n"
		buffer = tab(buffer, nest+1)
		buffer = buffer + "\"index\": " + strconv.Itoa(result.index) + ",\n"

		if result.resultType == "container" {
			buffer = tab(buffer, nest+1)
			buffer = buffer + "\"children\": [\n"
			var children = result.children
			for _, child := range children {
				buffer = loop(child, buffer, nest+2)
			}
			buffer = tab(buffer, nest+1)
			buffer = buffer + "],\n"
		} else {
			buffer = tab(buffer, nest+1)
			buffer = buffer + "\"value\": \"" + result.value + "\",\n"

			var expected = result.expected
			buffer = tab(buffer, nest+1)
			buffer = buffer + "\"expected\": [\n"
			for _, elm := range expected {
				buffer = tab(buffer, nest+2)
				buffer = buffer + "\"" + elm + "\",\n"
			}

			buffer = tab(buffer, nest+1)
			buffer = buffer + "],\n"
		}

		buffer = tab(buffer, nest)
		buffer = buffer + "}"

		if nest != 0 {
			buffer = buffer + ","
		}
		buffer = buffer + "\n"

		return buffer
	}

	return loop(result, "", 0)
}
