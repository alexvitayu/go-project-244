package compare

import (
	"fmt"
	"sort"
)

type Diff struct {
	IsNode   bool
	Path     string
	Message  string
	OldValue any
	NewValue any
}

func Compare(val1, val2 interface{}, basePath string) []Diff {
	var diffs []Diff
	obj1, isObj1 := IsObject(val1)
	obj2, isObj2 := IsObject(val2)
	arr1, isArr1 := IsArray(val1)
	arr2, isArr2 := IsArray(val2)
	switch {
	case isObj1 && isObj2:
		df := compareObjects(obj1, obj2, basePath)
		diffs = append(diffs, df...)
	case isArr1 && isArr2:
		df := compareArrays(arr1, arr2, basePath)
		diffs = append(diffs, df...)
	default:
		diffs = append(diffs, Diff{
			Path:     basePath,
			Message:  "mismatched types",
			OldValue: val1,
			NewValue: val2,
		})
	}
	return diffs
}

func compareArrays(arr1, arr2 []any, basePath string) []Diff {
	var diffs []Diff
	lenArr1 := len(arr1)
	lenArr2 := len(arr2)
	var minLen int
	if lenArr1 > lenArr2 {
		minLen = lenArr2
	} else {
		minLen = lenArr1
	}
	// Случай, пока количество элеменов в слайсе равно
	for i := 0; i < minLen; i++ {
		val1 := arr1[i]
		val2 := arr2[i]
		currentPath := basePath + fmt.Sprintf("[%d]", i)
		// раз оба значения val1 и val2 присутствуют в слайсе, то вызовем рекурсию, чтобы посмотреть, какие эти значения
		childDiff := Compare(val1, val2, currentPath)
		diffs = append(diffs, childDiff...)
	}
	// случай, когда элементы в слайсе 2 (новом) закончились
	if lenArr1 > lenArr2 {
		for i := minLen; i < lenArr1; i++ {
			currentPath := basePath + fmt.Sprintf("[%d]", i)
			df := Diff{
				Path:     currentPath,
				Message:  " # Удалена",
				OldValue: arr1[i],
				NewValue: nil,
			}
			diffs = append(diffs, df)
		}
		// случай, когда элементы в слайсе 1 (старом) закончились
	} else if lenArr1 < lenArr2 {
		for i := minLen; i < lenArr2; i++ {
			currentPath := basePath + fmt.Sprintf("[%d]", i)
			df := Diff{
				Path:     currentPath,
				Message:  " # Добавлена",
				OldValue: nil,
				NewValue: arr2[i],
			}
			diffs = append(diffs, df)
		}
	}
	return diffs
}

func compareObjects(val1, val2 any, basePath string) []Diff {
	var diffs []Diff
	obj1, _ := IsObject(val1)
	obj2, _ := IsObject(val2)
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

		currentPath := basePath + "." + key
		// Случай, когда оба значения объекты
		o1, isObj1 := IsObject(value1)
		o2, isObj2 := IsObject(value2)
		if ok && yes && isObj1 && isObj2 {
			childDiff := compareObjects(o1, o2, currentPath) // рекурсивный шаг
			df := Diff{
				Path:     currentPath,
				IsNode:   true,
				Message:  " # Объекты",
				OldValue: nil,
				NewValue: childDiff,
			}
			diffs = append(diffs, df)
			repeatedKeys[key] = true
			continue
			// остальные случаи
		} else {
			df := processSimpleCases(ok, yes, currentPath, value1, value2)
			diffs = append(diffs, df...)
		}
	}
	return diffs
}

func processSimpleCases(ok, yes bool, path string, val1, val2 interface{}) []Diff {
	var diffs []Diff
	_, isSimple1 := SimpleType(val1)
	_, isSimple2 := SimpleType(val2)
	switch {
	case ok && yes && isSimple1 && isSimple2 && val1 == val2:
		diffs = append(diffs, Diff{Path: path, Message: " # Равны", OldValue: val1, NewValue: val2})
	case ok && yes:
		diffs = append(diffs, Diff{Path: path, Message: " # Старое значение", OldValue: val1, NewValue: nil})
		diffs = append(diffs, Diff{Path: path, Message: " # Новое значение", OldValue: nil, NewValue: val2})
	case !ok && yes:
		diffs = append(diffs, Diff{Path: path, Message: " # Добавлена", OldValue: nil, NewValue: val2})
	case ok && !yes:
		diffs = append(diffs, Diff{Path: path, Message: " # Удалена", OldValue: val1, NewValue: nil})
	default:
		diffs = append(diffs, Diff{Path: path, Message: "", OldValue: val1, NewValue: val2})
	}
	return diffs
}

func IsObject(value interface{}) (map[string]any, bool) {
	v, ok := value.(map[string]interface{})
	return v, ok
}

func IsArray(value interface{}) ([]any, bool) {
	if value == nil {
		return nil, false
	}
	v, ok := value.([]interface{})
	return v, ok
}

func SimpleType(value interface{}) (any, bool) {
	switch value.(type) {
	case string, int, int8, int16, int32, int64, float32, float64, bool:
		return value, true
	default:
		return value, false
	}
}
