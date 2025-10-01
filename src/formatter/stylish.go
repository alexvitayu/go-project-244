package formatter

import (
	"code/src/compare"
	"fmt"
	"slices"
	"strings"
)

func FormatDiff(diff []compare.Diff) string {
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
	sortedKeys := normalizeKeys(m.DifTest)
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
			b.WriteString(SetIndent(m, depth) + "}" + "\n")
		case map[string]interface{}:
			b.WriteString(key + ": {" + "\n")
			b.WriteString(PrintObject(v, depth))
			b.WriteString(SetIndent(m, depth) + "}" + "\n")
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
func PrintObject(m map[string]interface{}, depth int) string {
	b := strings.Builder{}
	sortedKeys := normalizeKeys(m)
	for _, key := range sortedKeys {
		v := m[key] // Получаем значение по отсортированному ключу
		tempDiff := compare.Diff{
			DifTest: map[string]any{key: v},
		}
		b.WriteString(PrintMap(tempDiff, depth+1)) // Рекурсивный шаг
	}
	return b.String()
}

// normalizeKeys сортирует ключи в алфавитном порядке
func normalizeKeys(value map[string]any) []string {
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
		return diff.Type
	case " # Удалена", " # Старое значение":
		return diff.Type
	case "":
		return diff.Type
	default:
		return ".."
	}
}

// SetIndent устанавливает отступы для древовидной структуры
func SetIndent(diff compare.Diff, depth int) string {
	var indent string
	baseIndent := strings.Repeat("....", depth+1)
	shiftLeft := ".."
	//if strings.Contains(b.String(), "}") {
	//
	//	return baseIndent
	//}

	if diff.Type == "+ " || diff.Type == "- " || diff.Type == "  " {
		indent = strings.TrimSuffix(baseIndent, shiftLeft)
		return indent
	}
	return baseIndent
}

// финальная
//func FormatDiff(diff []compare.Diff) string {
//	b := strings.Builder{}
//	b.WriteString("{" + "\n")
//	for _, d := range diff {
//		str := PrintMap(d, 0)
//		b.WriteString(str)
//	}
//	b.WriteString("}")
//	return b.String()
//}
//func PrintMap(m compare.Diff, depth int) string {
//	b := strings.Builder{}
//	// Проведём нормализацию ключей для сохранения алфавитного порядка
//	sortedKeys := normalizeKeys(m.DifTest)
//	for _, key := range sortedKeys {
//		value := m.DifTest[key] // Получаем значение по отсортированному ключу
//
//		// Устанавливаем отступы для древовидной структуры
//		b.WriteString(SetIndent(m, depth))
//		// Устанавливаем префиксы +, - или " "
//		b.WriteString(SetPrefix(m))
//
//		switch v := value.(type) {
//		case []compare.Diff:
//			b.WriteString(key + ": {" + "\n")
//			b.WriteString(PrintArray(m, key, v, depth))
//			b.WriteString(SetIndent(m, depth) + "}" + "\n")
//		case map[string]interface{}:
//			b.WriteString(key + ": {" + "\n")
//			b.WriteString(PrintObject(v, depth))
//			b.WriteString(SetIndent(m, depth) + "..}" + "\n")
//		default:
//			b.WriteString(key + ": " + fmt.Sprint(value) + "\n")
//		}
//	}
//	return b.String()
//}
//
//// PrintArray реализует рекурсию по значениям, если они слайсы
//func PrintArray(m compare.Diff, key string, value []compare.Diff, depth int) string {
//	b := strings.Builder{}
//	for _, v := range value {
//		b.WriteString(PrintMap(v, depth+1)) // Рекурсивный шаг
//	}
//	return b.String()
//}
//
//// PrintObject реализует рекурсию по значениям, если они map-ы
//func PrintObject(m map[string]interface{}, depth int) string {
//	b := strings.Builder{}
//	sortedKeys := normalizeKeys(m)
//	for _, key := range sortedKeys {
//		v := m[key] // Получаем значение по отсортированному ключу
//		tempDiff := compare.Diff{
//			DifTest: map[string]any{key: v},
//		}
//		b.WriteString(PrintMap(tempDiff, depth+1)) // Рекурсивный шаг
//	}
//	return b.String()
//}
//
//// normalizeKeys сортирует ключи в алфавитном порядке
//func normalizeKeys(value map[string]any) []string {
//	var sortedKeys []string
//	for k := range value {
//		sortedKeys = append(sortedKeys, k)
//	}
//	slices.Sort(sortedKeys)
//	return sortedKeys
//}
//
//// SetPrefix устанавливает префиксы +, - или " "
//func SetPrefix(diff compare.Diff) string {
//	switch diff.Message {
//	case " # Добавлена", " # Новое значение":
//		return diff.Type
//	case " # Удалена", " # Старое значение":
//		return diff.Type
//	case "":
//		return diff.Type
//	default:
//		return ".."
//	}
//}
//
//// SetIndent устанавливает отступы для древовидной структуры
//func SetIndent(diff compare.Diff, depth int) string {
//	var indent string
//	baseIndent := strings.Repeat("....", depth+1)
//	shiftLeft := ".."
//
//	if diff.Type == "+ " || diff.Type == "- " || diff.Type == "  " {
//		indent = strings.TrimSuffix(baseIndent, shiftLeft)
//		return indent
//	}
//	return baseIndent
//}

// Супер, осталось только разобраться с отступом
//func FormatDiff(diff []compare.Diff) string {
//	b := strings.Builder{}
//	b.WriteString("{" + "\n")
//	for _, d := range diff {
//		str := PrintMap(d, 0)
//		b.WriteString(str)
//	}
//	b.WriteString("}")
//	return b.String()
//}
//func PrintMap(m compare.Diff, depth int) string {
//	b := strings.Builder{}
//
//	// Проведём нормализацию ключей для сохранения алфавитного порядка
//	sortedKeys := normalizeKeys(m.DifTest)
//	for _, key := range sortedKeys {
//		value := m.DifTest[key] // Получаем значение по отсортированному ключу
//
//		// Устанавливаем отступы для древовидной структуры
//		b.WriteString(SetIndent(m, depth))
//		// Устанавливаем префиксы +, - или " "
//		b.WriteString(SetPrefix(m))
//
//		switch v := value.(type) {
//		case []compare.Diff:
//			b.WriteString(key + ": {" + "\n")
//			b.WriteString(PrintArray(m, key, v, depth))
//			b.WriteString(SetIndent(m, depth) + "  }" + "\n")
//		case map[string]interface{}:
//			b.WriteString(key + ": {" + "\n")
//			b.WriteString(PrintObject(v, depth))
//			b.WriteString(SetIndent(m, depth) + "  }" + "\n")
//		default:
//			b.WriteString(key + ": " + fmt.Sprint(value) + "\n")
//		}
//	}
//	return b.String()
//}
//
//// PrintArray реализует рекурсию по значениям, если они слайсы
//func PrintArray(m compare.Diff, key string, value []compare.Diff, depth int) string {
//	b := strings.Builder{}
//	for _, v := range value {
//		b.WriteString(PrintMap(v, depth+1)) // Рекурсивный шаг
//	}
//	return b.String()
//}
//
//// PrintObject реализует рекурсию по значениям, если они map-ы
//func PrintObject(m map[string]interface{}, depth int) string {
//	b := strings.Builder{}
//	sortedKeys := normalizeKeys(m)
//	for _, key := range sortedKeys {
//		v := m[key] // Получаем значение по отсортированному ключу
//		tempDiff := compare.Diff{
//			DifTest: map[string]any{key: v},
//		}
//		b.WriteString(PrintMap(tempDiff, depth+1)) // Рекурсивный шаг
//	}
//	return b.String()
//}
//
//// normalizeKeys сортирует ключи в алфавитном порядке
//func normalizeKeys(value map[string]any) []string {
//	var sortedKeys []string
//	for k := range value {
//		sortedKeys = append(sortedKeys, k)
//	}
//	slices.Sort(sortedKeys)
//	return sortedKeys
//}
//
//// SetPrefix устанавливает префиксы +, - или " "
//func SetPrefix(diff compare.Diff) string {
//	switch diff.Message {
//	case " # Добавлена", " # Новое значение":
//		return "+ "
//	case " # Удалена", " # Старое значение":
//		return "- "
//	default:
//		return "  "
//	}
//}
//
//// SetIndent устанавливает отступы для древовидной структуры
//func SetIndent(diff compare.Diff, depth int) string {
//	var indent string
//	baseIndent := strings.Repeat("....", depth+1)
//	shiftLeft := ".."
//	indent = strings.TrimSuffix(baseIndent, shiftLeft)
//	return indent
//}

// лучше чем ОК
//func PrintMap(m compare.Diff, depth int) string {
//	b := strings.Builder{}
//	sortedKeys := normalizeKeys(m)
//
//	for _, key := range sortedKeys {
//		value := m.DifTest[key] // Получаем значение по отсортированному ключу
//
//		// Устанавливаем отступы для древовидной структуры
//		b.WriteString(SetIndent(m, depth))
//		// Устанавливаем префиксы +, - или " "
//		b.WriteString(SetPrefix(m))
//
//		switch v := value.(type) {
//		case []compare.Diff:
//			b.WriteString(PrintArray(m, key, v, depth))
//		case map[string]interface{}:
//			b.WriteString(PrintObject(m, key, v, depth))
//		default:
//			b.WriteString(key + ": " + fmt.Sprint(value) + "\n")
//		}
//	}
//	return b.String()
//}
//
//func PrintArray(m compare.Diff, key string, value []compare.Diff, depth int) string {
//	b := strings.Builder{}
//	b.WriteString(key + ": {" + "\n")
//	for _, v := range value {
//		b.WriteString(PrintMap(v, depth+1))
//	}
//	b.WriteString(SetIndent(m, depth) + "  }" + "\n")
//	return b.String()
//}
//
//func PrintObject(m compare.Diff, key string, value map[string]interface{}, depth int) string {
//	b := strings.Builder{}
//	b.WriteString(key + ": {" + "\n")
//	for k, v := range value {
//		tempDiff := compare.Diff{
//			DifTest: map[string]any{k: v},
//		}
//		b.WriteString(PrintMap(tempDiff, depth+1))
//		b.WriteString(SetIndent(m, depth) + "  }" + "\n")
//	}
//	return b.String()
//}
//
//// normalizeKeys сортирует ключи в алфавитном порядке
//func normalizeKeys(m compare.Diff) []string {
//	var sortedKeys []string
//	for k := range m.DifTest {
//		sortedKeys = append(sortedKeys, k)
//	}
//	slices.Sort(sortedKeys)
//	return sortedKeys
//}
//
//// SetPrefix устанавливает префиксы +, - или " "
//func SetPrefix(diff compare.Diff) string {
//	switch diff.Message {
//	case " # Добавлена", " # Новое значение":
//		return "+ "
//	case " # Удалена", " # Старое значение":
//		return "- "
//	default:
//		return "  "
//	}
//}
//
//// SetIndent устанавливает отступы для древовидной структуры
//func SetIndent(diff compare.Diff, depth int) string {
//	var indent string
//	baseIndent := strings.Repeat("....", depth+1)
//	shiftLeft := ".."
//	indent = strings.TrimSuffix(baseIndent, shiftLeft)
//	return indent
//}

// Почти ОК
//func FormatDiff(diff []compare.Diff) string {
//	b := strings.Builder{}
//	b.WriteString("{" + "\n")
//	for _, d := range diff {
//		str := PrintMap(d, 1)
//		b.WriteString(str)
//	}
//	b.WriteString("}")
//	return b.String()
//}
//
//func PrintMap(m compare.Diff, depth int) string {
//	b := strings.Builder{}
//	for key, value := range m.DifTest {
//		// Устанавливаем отступы для древовидной структуры
//		b.WriteString(SetIndent(m, depth))
//		// Устанавливаем префиксы +, - или " "
//		b.WriteString(SetPrefix(m))
//		switch v := value.(type) {
//		case []compare.Diff:
//			b.WriteString(PrintArray(m, key, v, depth))
//		case map[string]interface{}:
//			b.WriteString(PrintObject(m, key, v, depth))
//		default:
//			b.WriteString(key + ": " + fmt.Sprint(value) + "\n")
//		}
//	}
//	return b.String()
//}
//
//func PrintArray(m compare.Diff, key string, value []compare.Diff, depth int) string {
//	b := strings.Builder{}
//	b.WriteString(key + ": {" + "\n")
//	for _, v := range value {
//		b.WriteString(PrintMap(v, depth+1))
//	}
//	b.WriteString(SetIndent(m, depth) + "}" + "\n")
//	return b.String()
//}
//
//func PrintObject(m compare.Diff, key string, value map[string]interface{}, depth int) string {
//	b := strings.Builder{}
//	b.WriteString(key + ": {" + "\n")
//	for k, v := range value {
//		tempDiff := compare.Diff{
//			DifTest: map[string]any{k: v},
//		}
//		b.WriteString(PrintMap(tempDiff, depth+1))
//		b.WriteString(SetIndent(m, depth) + "  }" + "\n")
//	}
//	return b.String()
//}
//
//// SetPrefix устанавливает префиксы +, - или " "
//func SetPrefix(diff compare.Diff) string {
//	switch diff.Message {
//	case " # Добавлена", " # Новое значение":
//		return "+ "
//	case " # Удалена", " # Старое значение":
//		return "- "
//	default:
//		return ""
//	}
//}
//
//// SetIndent устанавливает отступы для древовидной структуры
//func SetIndent(diff compare.Diff, depth int) string {
//	var indent string
//	baseIndent := strings.Repeat("....", depth)
//	shiftLeft := ".."
//	if depth > 0 {
//		if diff.ShiftLeft {
//			indent = strings.TrimSuffix(baseIndent, shiftLeft)
//			return indent
//		} else {
//			return baseIndent
//		}
//	}
//	return ".."
//}

// второй рабочий вариант с упорядоченной архитектурой
//func FormatDiff(diff []compare.Diff) string {
//	b := strings.Builder{}
//	for _, d := range diff {
//		str := PrintMap(d)
//		b.WriteString(str)
//	}
//	return b.String()
//}
//
//func PrintMap(m compare.Diff) string {
//	b := strings.Builder{}
//	for key, value := range m.DifTest {
//		switch v := value.(type) {
//		case []compare.Diff:
//
//			b.WriteString(PrintArray(key, v))
//		case map[string]interface{}:
//
//			b.WriteString(PrintObject(key, v))
//		default:
//			b.WriteString(key + ":" + fmt.Sprint(value) + "\n")
//		}
//	}
//	return b.String()
//}
//
//func PrintArray(key string, value []compare.Diff) string {
//	b := strings.Builder{}
//	b.WriteString(key + "\n")
//	for _, v := range value {
//		b.WriteString(PrintMap(v))
//	}
//	return b.String()
//}
//
//func PrintObject(key string, value map[string]interface{}) string {
//	b := strings.Builder{}
//	b.WriteString(key + "\n")
//	for k, v := range value {
//		tempDiff := compare.Diff{
//			DifTest: map[string]any{k: v},
//		}
//		b.WriteString(PrintMap(tempDiff))
//	}
//	return b.String()
//}

// первый рабочий вариант неупорядоченная архитектура
//func FormatDiff(diff []compare.Diff) string {
//	b := strings.Builder{}
//	for _, d := range diff {
//		str := PrintMap(d)
//		b.WriteString(str)
//	}
//	return b.String()
//}
//
//func PrintMap(m compare.Diff) string {
//	b := strings.Builder{}
//	for key, value := range m.DifTest {
//		switch v := value.(type) {
//		case []compare.Diff:
//			b.WriteString(key + "\n")
//			str := PrintArray(v)
//			b.WriteString(str + "\n")
//		case map[string]interface{}:
//			b.WriteString(key + "\n")
//			str := PrintObject(v, m)
//			b.WriteString(str + "\n")
//		default:
//			b.WriteString(key + ":" + fmt.Sprint(value) + m.Message + "\n")
//		}
//	}
//	return b.String()
//}
//
//func PrintArray(m []compare.Diff) string {
//	b := strings.Builder{}
//	for _, v := range m {
//		for key, value := range v.DifTest {
//			array, ok := value.([]compare.Diff)
//			if ok {
//				b.WriteString(key + "\n")
//				str := PrintArray(array)
//				b.WriteString(str + "\n")
//			} else if compare.IsObject(value) {
//				b.WriteString(key + "\n")
//				obj := value.(map[string]interface{})
//				str := PrintObject(obj, v)
//				b.WriteString(str + "\n")
//			} else {
//				b.WriteString(key + ":" + fmt.Sprint(value) + v.Message + "\n")
//			}
//		}
//	}
//	return b.String()
//}
//
//func PrintObject(m map[string]interface{}, o compare.Diff) string {
//	b := strings.Builder{}
//	for key, value := range m {
//		if compare.IsObject(value) {
//			b.WriteString(key + "\n")
//			obj := value.(map[string]interface{})
//			str := PrintObject(obj, o)
//			b.WriteString(str + "\n")
//		} else if compare.IsArray(value) {
//			b.WriteString(key + "\n")
//			obj := value.([]compare.Diff)
//			str := PrintArray(obj)
//			b.WriteString(str + "\n")
//		} else {
//			b.WriteString(key + ":" + fmt.Sprint(value) + o.Message + "\n")
//		}
//	}
//	return b.String()
//}
