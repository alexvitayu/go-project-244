package formatters

import (
	"code/src/compare"
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type Property struct {
	path    string
	message string
	value1  any
	value2  any
}

type PlainFormatter struct{}

func NewPlainFormatter() *PlainFormatter {
	return &PlainFormatter{}
}

func (st *PlainFormatter) FormatDiff(diffs []compare.Diff) (string, error) {
	b := strings.Builder{}
	b.WriteString("\n")
	keys := processKeys(diffs)
	properties, _ := process(diffs, keys, 0)
	for _, p := range properties {
		b.WriteString(printPlain(p))
	}
	return b.String(), nil
}

func printPlain(p Property) string {
	b := strings.Builder{}
	var val1 any
	var val2 any
	if !compare.SimpleType(p.value1) && p.value1 != nil {
		val1 = printValue("[complex value]")
	} else if p.value1 == nil {
		val1 = nil
	} else {
		val1 = printValue(p.value1)
	}
	if !compare.SimpleType(p.value2) && p.value2 != nil {
		val2 = printValue("[complex value]")
	} else if p.value2 == nil {
		val2 = nil
	} else {
		val2 = printValue(p.value2)
	}
	switch p.message {
	case " # Добавлена":
		b.WriteString("Property" + color.YellowString(fmt.Sprintf(" '%v' ", strings.TrimPrefix(p.path, "."))) + "was added with value: " + fmt.Sprint(val1) + "\n")
	case " # Новое значение", " # Старое значение":
		b.WriteString("Property" + color.YellowString(fmt.Sprintf(" '%v' ", strings.TrimPrefix(p.path, "."))) + "was updated. " + fmt.Sprintf("From %v to %v", val1, val2) + "\n")
	case " # Удалена":
		b.WriteString("Property" + color.YellowString(fmt.Sprintf(" '%v' ", strings.TrimPrefix(p.path, "."))) + "was removed" + "\n")
	}
	return b.String()
}

func printValue(val any) string {
	if val == "[complex value]" {
		val = color.RedString("[") + "complex value]"
		return fmt.Sprint(val)
	}
	switch val.(type) {
	case bool:
		return fmt.Sprint(val)
	case nil:
		return fmt.Sprint(nil)
	default:
		return color.YellowString(fmt.Sprintf("'%v'", val))
	}
}

func processKeys(diffs []compare.Diff) []string {
	var allKeys []string
	for _, diff := range diffs {
		for key, value := range diff.DifTest {
			allKeys = append(allKeys, key)
			_, ok := value.([]compare.Diff)
			if ok {
				keys := processKeys(value.([]compare.Diff))
				allKeys = append(allKeys, keys...)
			}
		}
	}
	return allKeys
}

func process(diffs []compare.Diff, keys []string, startIndex int) ([]Property, int) {
	properties := []Property{}
	index := startIndex
	for j := 0; j < len(diffs) && index < len(keys); j++ {
		currentDiff := diffs[j]
		currentValue := currentDiff.DifTest[keys[index]]
		switch v := currentValue.(type) {
		case []compare.Diff:
			// Обрабатываем элемент верхнего уровня
			pr := Property{
				path:    currentDiff.Path,
				message: currentDiff.Message,
				value1:  nil,
				value2:  nil,
			}
			properties = append(properties, pr)
			childProp, newIndex := process(v, keys, index+1)
			properties = append(properties, childProp...)
			index = newIndex
		default:
			if j+1 < len(diffs) &&
				index+1 < len(keys) &&
				currentDiff.Message == " # Старое значение" &&
				diffs[j+1].Message == " # Новое значение" &&
				diffs[j].Path == diffs[j+1].Path {
				// Создаем Property со старым и новым значением
				nextDiff := diffs[j+1]
				nextValue := nextDiff.DifTest[keys[index+1]]
				prop := Property{
					path:    currentDiff.Path,
					message: currentDiff.Message,
					value1:  currentValue, // старое значение
					value2:  nextValue,    // новое значение
				}
				properties = append(properties, prop)
				index += 2 //пропускаем пару ключей
				j++        // Пропускаем следующий элемент, так как он уже обработан
			} else {
				// Это случай для одиночного значения
				prop := Property{
					path:    currentDiff.Path,
					message: currentDiff.Message,
					value1:  currentValue,
					value2:  nil,
				}
				properties = append(properties, prop)
				index++
			}
		}
	}
	return properties, index
}
