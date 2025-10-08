package formatters

import (
	"code/src/compare"

	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type HexletJson struct {
	Key      string       `json:"key"`
	Type     string       `json:"type"`
	Value1   interface{}  `json:"value1,omitempty"`
	Value2   interface{}  `json:"value2,omitempty"`
	Children []HexletJson `json:"children,omitempty"`
}
type JsonFormater struct{}

func NewJsonFormatter() *JsonFormater {
	return &JsonFormater{}
}

func (j *JsonFormater) FormatDiff(diffs []compare.Diff) (string, error) {
	jDiffs := processJson(diffs)
	jSn, err := printJson(jDiffs)
	if err != nil {
		return "", fmt.Errorf("printJson: %w", err)
	}
	return jSn, nil
}

func printJson(jDiffs []HexletJson) (string, error) {
	jD := HexletJson{
		Key:      "",
		Type:     "root",
		Value1:   nil,
		Value2:   nil,
		Children: jDiffs,
	}
	bytes, err := json.MarshalIndent(jD, "", "  ")
	if err != nil {
		return "", errors.New("не удалось преобразовать в json")
	}
	return string(bytes), nil
}

func processJson(diffs []compare.Diff) []HexletJson {
	var jDiffs []HexletJson
	for i := 0; i < len(diffs); i++ {
		currentDiff := diffs[i]
		parts := strings.Split(currentDiff.Path, ".")
		currentKey := parts[len(parts)-1]
		currentVal := check(currentDiff.OldValue, currentDiff.NewValue)
		switch v := currentVal.(type) {
		case []compare.Diff:
			childJDiff := processJson(v)
			jD := HexletJson{
				Key:      currentKey,
				Type:     addMessage(currentDiff),
				Value1:   nil,
				Value2:   nil,
				Children: childJDiff,
			}
			jDiffs = append(jDiffs, jD)

		default:
			if i+1 < len(diffs) &&
				// Обрабатываем случай с новым и старым значением
				currentDiff.Message == " # Старое значение" &&
				diffs[i+1].Message == " # Новое значение" &&
				currentDiff.Path == diffs[i+1].Path {
				nextDiff := diffs[i+1]
				//nextKey := currentKey
				jD := HexletJson{
					Key:    currentKey,
					Type:   addMessage(currentDiff),
					Value1: currentDiff.OldValue,
					Value2: nextDiff.NewValue,
				}
				jDiffs = append(jDiffs, jD)
				i++ // пропускаем один дифф с новым значением, т.к. он уже обработан
			} else if currentDiff.Message == " # Равны" {
				jD := HexletJson{
					Key:    currentKey,
					Type:   addMessage(currentDiff),
					Value1: currentDiff.OldValue,
				}
				jDiffs = append(jDiffs, jD)
			} else {
				// Это случай для одиночного значения
				jD := HexletJson{
					Key:    currentKey,
					Type:   addMessage(currentDiff),
					Value1: currentDiff.OldValue,
					Value2: currentDiff.NewValue,
				}
				jDiffs = append(jDiffs, jD)
			}
		}
	}
	return jDiffs
}

func addMessage(diff compare.Diff) string {
	switch diff.Message {
	case " # Объекты":
		return "nested"
	case " # Добавлена":
		return "added"
	case " # Удалена":
		return "deleted"
	case " # Старое значение", " # Новое значение":
		return "changed"
	case " # Равны":
		return "unchanged"
	default:
		return ""
	}
}
