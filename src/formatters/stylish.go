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

func (st *Stylish) FormatDiff(diff []compare.Diff) (string, error) {
	b := strings.Builder{}
	b.WriteString("{" + "\n")
	b.WriteString(printStylish(diff, 0))
	b.WriteString("}")
	return b.String(), nil
}

func printStylish(df []compare.Diff, depth int) string {
	b := strings.Builder{}
	for _, d := range df {
		currentPath := d.Path
		parts := strings.Split(currentPath, ".")
		key := parts[len(parts)-1]
		value := check(d.OldValue, d.NewValue)
		// Устанавливаем отступы для древовидной структуры
		b.WriteString(setIndent(d, depth))
		// Устанавливаем префиксы +, - или " "
		b.WriteString(setPrefix(d))
		switch v := value.(type) {
		case []compare.Diff:
			b.WriteString(key + ": {" + "\n")
			b.WriteString(printStylish(v, depth+1))
			b.WriteString(adjustIndent(depth) + "}" + "\n")
		case map[string]interface{}:
			var diffs []compare.Diff
			b.WriteString(key + ": {" + "\n")
			sortedKeys := normalizeKeys(v)
			for _, key := range sortedKeys {
				val := v[key]
				path := currentPath + "." + key
				diffs = append(diffs, compare.Diff{
					Path:     path,
					OldValue: nil,
					NewValue: val,
				})
			}
			b.WriteString(printStylish(diffs, depth+1))
			b.WriteString(adjustIndent(depth) + "}" + "\n")
		default:
			b.WriteString(key + ": " + formatValue(extractValue(d)) + "\n")
		}
	}
	return b.String()
}

// check проверяет, oldValue и newValue в диффе и возвращает не nil-овое значение
func check(val1, val2 any) any {
	if val1 == nil {
		return val2
	} else if val2 == nil {
		return val1
	} else {
		return val2
	}
}

// normalizeKeys сортирует ключи из мапы в алфавитном порядке
func normalizeKeys(value map[string]any) []string {
	var sortedKeys []string
	for k := range value {
		sortedKeys = append(sortedKeys, k)
	}
	slices.Sort(sortedKeys)
	return sortedKeys
}

// extractValue заполняет поля структуры oldValue и newValue в зависимости от сообщения в диффе
func extractValue(df compare.Diff) any {
	switch df.Message {
	case " # Равны":
		return df.OldValue
	case " # Удалена", " # Старое значение":
		return df.OldValue
	case " # Добавлена", " # Новое значение":
		return df.NewValue
	default:
		return df.NewValue
	}
}

func formatValue(val any) string {
	switch val.(type) {
	case nil:
		return "null"
	default:
		return fmt.Sprint(val)
	}
}

// setPrefix устанавливает префиксы +, - или " "
func setPrefix(diff compare.Diff) string {
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

// setIndent устанавливает отступы для древовидной структуры
func setIndent(diff compare.Diff, depth int) string {
	var indent string
	baseIndent := strings.Repeat("    ", depth+1)
	shiftLeft := "  "
	if diff.Message == " # Добавлена" ||
		diff.Message == " # Старое значение" ||
		diff.Message == " # Новое значение" ||
		diff.Message == " # Удалена" ||
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

// adjustIndent регулирует устанавливает отступы для закрывающейся скобки
func adjustIndent(depth int) string {
	baseIndent := strings.Repeat("    ", depth+1)
	return baseIndent
}
