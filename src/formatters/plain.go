package formatters

import (
	"code/src/compare"
	"fmt"
	"strings"
)

type PlainFormatter struct{}

func NewPlainFormatter() *PlainFormatter {
	return &PlainFormatter{}
}

func (st *PlainFormatter) FormatDiff(diffs []compare.Diff) string {
	b := strings.Builder{}
	for _, df := range diffs {
		b.WriteString(Plain(df, 0))
	}
	return b.String()
}

func PrintPlain(diff compare.Diff, value1, value2 any) string {
	b := strings.Builder{}
	switch diff.Message {
	case " # Добавлена":
		b.WriteString("Property" + fmt.Sprintf("%v"+"", diff.Path) + "was added with value: " + fmt.Sprint(value2) + "\n")
	case " # Новое значение", " # Старое значение":
		b.WriteString("Property" + fmt.Sprintf("%v"+"", diff.Path) + "was updated. " + fmt.Sprintf("from '%v' to '%v", value1, value2) + "\n")
	case " # Удалена":
		b.WriteString("Property" + fmt.Sprintf("%v"+"", diff.Path) + "was removed with value: " + fmt.Sprint(value1) + "\n")
	}
	return b.String()
}

func Plain(diff compare.Diff, depth int) string {
	b := strings.Builder{}
	//keys := NormalizeKeys(diff.DifTest)
	//
	//for _, key := range keys {
	//
	//}

	return b.String()
}
