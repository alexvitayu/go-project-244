package formatters

import (
	"code/src/compare"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type JsonDiff struct {
	Path     string `json:"path"`
	Message  string `json:"message"`
	OldValue any    `json:"oldValue"`
	NewValue any    `json:"newValue"`
}

type JsonFormatter struct{}

func NewJsonFormatter() *JsonFormatter {
	return &JsonFormatter{}
}

func (j *JsonFormatter) FormatDiff(diffs []compare.Diff) (string, error) {
	jDiffs := processJson(diffs)
	jSn, err := printJson(jDiffs)
	if err != nil {
		return "", fmt.Errorf("printJson: %w", err)
	}
	return jSn, nil
}

func printJson(jDiffs []JsonDiff) (string, error) {
	bytes, err := json.MarshalIndent(jDiffs, "", "  ")
	if err != nil {
		return "", errors.New("не удалось преобразовать в json")
	}
	return string(bytes), nil
}

func processJson(diffs []compare.Diff) []JsonDiff {
	var jDiffs []JsonDiff
	for i := 0; i < len(diffs); i++ {
		currentDiff := diffs[i]
		parts := strings.Split(currentDiff.Path, ".")
		currentKey := parts[len(parts)-1]
		currentVal := currentDiff.DifTest[currentKey]
		switch v := currentVal.(type) {
		case []compare.Diff:
			childJDiff := processJson(v)
			jD := JsonDiff{
				Path:     currentDiff.Path,
				Message:  addMessage(currentDiff),
				OldValue: nil,
				NewValue: childJDiff,
			}
			jDiffs = append(jDiffs, jD)

		default:
			if currentDiff.Message == " # Равны" {
				continue
			}
			if i+1 < len(diffs) &&
				currentDiff.Message == " # Старое значение" &&
				diffs[i+1].Message == " # Новое значение" &&
				currentDiff.Path == diffs[i+1].Path {
				// Обрабатываем случай с новым и старым значением
				nextDiff := diffs[i+1]
				nextKey := currentKey
				jD := JsonDiff{
					Path:     currentDiff.Path,
					Message:  addMessage(currentDiff),
					OldValue: oldValue(currentDiff, currentKey),
					NewValue: newValue(nextDiff, nextKey),
				}
				jDiffs = append(jDiffs, jD)
				i++ // пропускаем один дифф с новым значением, т.к. он уже обработан
			} else {
				// Это случай для одиночного значения
				jD := JsonDiff{
					Path:     currentDiff.Path,
					Message:  addMessage(currentDiff),
					OldValue: oldValue(currentDiff, currentKey),
					NewValue: newValue(currentDiff, currentKey),
				}
				jDiffs = append(jDiffs, jD)
			}
		}
	}
	return jDiffs
}

func processVal(val any) any {
	if val == nil {
		return nil
	}
	return val
}

func oldValue(diff compare.Diff, key string) any {
	val := diff.DifTest[key]
	switch diff.Message {
	case " # Добавлена", " # Новое значение":
		return nil
	case " # Удалена", " # Старое значение":
		return processVal(val)
	default:
		return nil
	}
}
func newValue(diff compare.Diff, key string) any {
	val := diff.DifTest[key]
	switch diff.Message {
	case " # Добавлена", " # Новое значение":
		return processVal(val)
	case " # Удалена", " # Старое значение":
		return nil
	default:
		return processVal(val)
	}
}

func addMessage(diff compare.Diff) string {
	switch diff.Message {
	case " # Объекты":
		return "objects"
	case " # Добавлена":
		return "added"
	case " # Удалена":
		return "deleted"
	case " # Старое значение", " # Новое значение":
		return "modified"
	default:
		return ""
	}
}
