package formatter

import (
	"code/src/compare"
	"fmt"
	"strings"
)

func FormatDiff(diff []compare.Diff) string {
	b := strings.Builder{}
	for _, d := range diff {
		str := PrintMap(d)
		b.WriteString(str + "\n")
	}
	return b.String()
}

func PrintMap(m compare.Diff) string {
	b := strings.Builder{}
	for key, value := range m.DifTest {
		switch v := value.(type) {
		case []compare.Diff:
			b.WriteString(key + "\n")
			str := PrintArray(v)
			b.WriteString(str + "\n")
		case map[string]interface{}:
			b.WriteString(key + "\n")
			str := PrintObject(v, m)
			b.WriteString(str + "\n")
		default:
			b.WriteString(key + ":" + fmt.Sprint(value) + m.Message + "\n")
		}
	}
	return b.String()
}

func PrintArray(m []compare.Diff) string {
	b := strings.Builder{}
	for _, v := range m {
		for key, value := range v.DifTest {
			array, ok := value.([]compare.Diff)
			if ok {
				b.WriteString(key + "\n")
				str := PrintArray(array)
				b.WriteString(str + "\n")
			} else if compare.IsObject(value) {
				b.WriteString(key + "\n")
				obj := value.(map[string]interface{})
				str := PrintObject(obj, v)
				b.WriteString(str + "\n")
			} else {
				b.WriteString(key + ":" + fmt.Sprint(value) + v.Message + "\n")
			}
		}
	}
	return b.String()
}

func PrintObject(m map[string]interface{}, o compare.Diff) string {
	b := strings.Builder{}
	for key, value := range m {
		if compare.IsObject(value) {
			b.WriteString(key + "\n")
			obj := value.(map[string]interface{})
			str := PrintObject(obj, o)
			b.WriteString(str + "\n")
		} else if compare.IsArray(value) {
			b.WriteString(key + "\n")
			obj := value.([]compare.Diff)
			str := PrintArray(obj)
			b.WriteString(str + "\n")
		} else {
			b.WriteString(key + ":" + fmt.Sprint(value) + o.Message + "\n")
		}
	}
	return b.String()
}
