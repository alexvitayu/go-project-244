package formatters

import (
	"code/src/compare"
	"fmt"
	"slices"
	"strings"
)

type Stylish struct{}

func NewStylish() *Stylish {
	return &Stylish{}
}

func (st *Stylish) FormatDiff(diff []compare.Diff) string {
	b := strings.Builder{}
	b.WriteString("{" + "\n")
	for _, d := range diff {
		str := PrintMap(d, 0)
		b.WriteString(str)
	}
	b.WriteString("}")
	return b.String()
}
func PrintMap(m compare.Diff, depth int) string {
	b := strings.Builder{}
	// Проведём нормализацию ключей для сохранения алфавитного порядка
	sortedKeys := NormalizeKeys(m.DifTest)
	for _, key := range sortedKeys {
		value := m.DifTest[key] // Получаем значение по отсортированному ключу

		// Устанавливаем отступы для древовидной структуры
		b.WriteString(SetIndent(m, depth))
		// Устанавливаем префиксы +, - или " "
		b.WriteString(SetPrefix(m))
		switch v := value.(type) {
		case []compare.Diff:
			b.WriteString(key + ": {" + "\n")
			b.WriteString(PrintArray(m, key, v, depth))
			b.WriteString(AdjustIndent(m, depth) + "}" + "\n")
		case map[string]interface{}:
			b.WriteString(key + ": {" + "\n")
			b.WriteString(PrintObject(m, v, depth))
			b.WriteString(AdjustIndent(m, depth) + "}" + "\n")
		default:
			b.WriteString(key + ": " + fmt.Sprint(value) + "\n")
		}
	}
	return b.String()
}

// PrintArray реализует рекурсию по значениям, если они слайсы
func PrintArray(m compare.Diff, key string, value []compare.Diff, depth int) string {
	b := strings.Builder{}
	for _, v := range value {
		b.WriteString(PrintMap(v, depth+1)) // Рекурсивный шаг
	}
	return b.String()
}

// PrintObject реализует рекурсию по значениям, если они map-ы
func PrintObject(a compare.Diff, m map[string]interface{}, depth int) string {
	b := strings.Builder{}
	sortedKeys := NormalizeKeys(m)
	for _, key := range sortedKeys {
		v := m[key] // Получаем значение по отсортированному ключу
		tempDiff := compare.Diff{
			DifTest: map[string]any{key: v},
		}
		b.WriteString(PrintMap(tempDiff, depth+1)) // Рекурсивный шаг
	}
	return b.String()
}

// NormalizeKeys сортирует ключи в алфавитном порядке
func NormalizeKeys(value map[string]any) []string {
	var sortedKeys []string
	for k := range value {
		sortedKeys = append(sortedKeys, k)
	}
	slices.Sort(sortedKeys)
	return sortedKeys
}

// SetPrefix устанавливает префиксы +, - или " "
func SetPrefix(diff compare.Diff) string {
	switch diff.Message {
	case " # Добавлена", " # Новое значение":
		return "+ "
	case " # Удалена", " # Старое значение":
		return "- "
	case " # Равны":
		return "  "
	case " # Объекты":
		return "  "
	default:
		return ""
	}
}

// SetIndent устанавливает отступы для древовидной структуры
func SetIndent(diff compare.Diff, depth int) string {
	var indent string
	baseIndent := strings.Repeat("    ", depth+1)
	shiftLeft := "  "
	if diff.Message == " # Добавлена" ||
		diff.Message == " # Новое значение" ||
		diff.Message == " # Удалена" ||
		diff.Message == " # Старое значение" ||
		diff.Message == " # Равны" ||
		diff.Message == " # Объекты" {
		indent = strings.TrimSuffix(baseIndent, shiftLeft)
		return indent
	} else if diff.Message == "" && !diff.IsNode {
		return baseIndent
	} else if diff.Message == " # Объекты" && diff.IsNode {
		return baseIndent
	}
	return ""
}

func AdjustIndent(diff compare.Diff, depth int) string {
	var indent string
	baseIndent := strings.Repeat("    ", depth+1)
	shiftLeft := "  "
	if diff.IsNode {
		return baseIndent
	} else if diff.Message == " # Добавлена" ||
		diff.Message == " # Новое значение" ||
		diff.Message == " # Удалена" ||
		diff.Message == " # Старое значение" ||
		diff.Message == " # Равны" ||
		diff.Message == " # Объекты" {
		indent = strings.TrimSuffix(baseIndent, shiftLeft)
		return indent
	}
	return baseIndent
}
