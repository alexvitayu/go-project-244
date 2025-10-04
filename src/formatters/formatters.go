package formatters

import "code/src/compare"

type Formater interface {
	FormatDiff([]compare.Diff) string
}

func Format(flag string) Formater {
	switch flag {
	case "plain":
		return NewPlainFormatter()
	case "json":
		return NewJsonFormatter()
	default:
		return NewStylish()
	}
}
