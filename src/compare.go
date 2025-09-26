package src

import (
	"fmt"
	"sort"
)

type Diff struct {
	Difference string
	Message    string
}

func Compare(obj1, obj2 map[string]interface{}, basePath string) []Diff {
	var diffs []Diff

	// Создаём множество ключей из обеих map
	uniqueKeys := make(map[string]struct{}) // Убираю дубликаты ключей
	for key := range obj1 {
		uniqueKeys[key] = struct{}{}
	}
	for key := range obj2 {
		uniqueKeys[key] = struct{}{}
	}
	commonKeys := make([]string, 0, len(uniqueKeys)) //Сразу аллоцирую память под известное количество элементов
	for key := range uniqueKeys {
		commonKeys = append(commonKeys, key)
	}
	// Сортируем ключи в алфавитном порядке
	sort.Slice(commonKeys, func(i, j int) bool {
		return commonKeys[i] < commonKeys[j]
	})
	repeatedKeys := make(map[string]bool) //для обработки повторяющихся ключей
	for _, key := range commonKeys {
		if repeatedKeys[key] {
			continue
		}

		value1, ok := obj1[key]
		value2, yes := obj2[key]
		currentPath := buildPath(basePath, key)

		if isObject(value2) && isObject(value1) {
			o1 := value1.(map[string]interface{})
			o2 := value2.(map[string]interface{})
			childDiff := Compare(o1, o2, currentPath) // рекурсивный шаг
			diffs = append(diffs, childDiff...)
			repeatedKeys[key] = true
			continue
		}
		switch {
		case !ok && yes:
			strData2 := fmt.Sprintf(" + %v: %v", key, value2)
			diffs = append(diffs, Diff{Difference: strData2, Message: "# Добавлена"})
		case ok && !yes:
			strData1 := fmt.Sprintf(" - %v: %v", key, value1)
			diffs = append(diffs, Diff{Difference: strData1, Message: "# Удалена"})

		default:
			isSimple1 := simpleType(value1)
			isSimple2 := simpleType(value2)
			// Оба простых типа и значения не равны
			if isSimple1 && isSimple2 && value1 != value2 {
				strData1 := fmt.Sprintf(" - %v: %v", key, value1)
				diffs = append(diffs, Diff{Difference: strData1, Message: "# Старое значение"})
				strData2 := fmt.Sprintf(" + %v: %v", key, value2)
				diffs = append(diffs, Diff{Difference: strData2, Message: "# Новое значение"})
				// Оба простых типа и значения равны
			} else if isSimple1 && isSimple2 && value1 == value2 {
				strData1 := fmt.Sprintf("   %v: %v", key, value1)
				diffs = append(diffs, Diff{Difference: strData1, Message: ""})
				// Если оба типа и не простые типы и не объекты
			} else {
				strData1 := fmt.Sprintf(" - %v: %v", key, value1)
				diffs = append(diffs, Diff{Difference: strData1, Message: "# Старое значение"})
				strData2 := fmt.Sprintf(" + %v: %v", key, value2)
				diffs = append(diffs, Diff{Difference: strData2, Message: "# Новое значение"})
			}
		}
	}
	return diffs
}

func buildPath(basePath, key string) string {
	if basePath == "" {
		return key
	}
	return basePath + "." + key
}

func simpleType(value interface{}) bool {
	switch value.(type) {
	case string, int, int8, int16, int32, int64, float32, float64, bool:
		return true
	default:
		return false
	}
}

func isArray(value interface{}) bool {
	_, ok := value.([]interface{})
	if ok {
		return true
	} else {
		return false
	}
}

func isObject(value interface{}) bool {
	_, ok := value.(map[string]interface{})
	if ok {
		return true
	} else {
		return false
	}
}
