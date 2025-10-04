package formatters

import (
	"code/src/compare"
)

type JsonFormatter struct{}

func NewJsonFormatter() *JsonFormatter {
	return &JsonFormatter{}
}

func (j *JsonFormatter) FormatDiff(diffs []compare.Diff) string {

	return ""
}
