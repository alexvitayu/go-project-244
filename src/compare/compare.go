package compare

import (
	"sort"
)

type Diff struct {
	IsNode  bool
	Path    string
	Message string
	DifTest map[string]any
}

func Compare(o1, o2 any, basePath string) []Diff {
	var diffs []Diff
	obj1, _ := o1.(map[string]interface{})
	obj2, _ := o2.(map[string]interface{})
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

		path := basePath + "." + key

		if ok && yes && IsObject(value1) && IsObject(value2) {
			node := make(map[string]any)

			o1 := value1.(map[string]interface{})
			o2 := value2.(map[string]interface{})
			childDiff := Compare(o1, o2, path) // рекурсивный шаг
			node[key] = childDiff
			dif := Diff{
				DifTest: map[string]any{
					key: childDiff,
				},
				Path:    path,
				IsNode:  true,
				Message: " # Объекты",
			}
			diffs = append(diffs, dif)
			repeatedKeys[key] = true
			continue
		}
		isSimple1 := SimpleType(value1)
		isSimple2 := SimpleType(value2)
		isObject1 := IsObject(value1)
		isObject2 := IsObject(value2)
		switch {
		case ok && yes && isSimple1 && isObject2:
			diffs = append(diffs, Diff{DifTest: map[string]any{key: value1}, Message: " # Старое значение", Path: path})
			diffs = append(diffs, Diff{DifTest: map[string]any{key: value2}, Message: " # Новое значение", Path: path, IsNode: true})
		case ok && yes && isObject1 && isSimple2:
			diffs = append(diffs, Diff{DifTest: map[string]any{key: value1}, Message: " # Старое значение", Path: path, IsNode: true})
			diffs = append(diffs, Diff{DifTest: map[string]any{key: value2}, Message: " # Новое значение", Path: path})
		case !ok && yes && isObject2:
			diffs = append(diffs, Diff{DifTest: map[string]any{key: value2}, Message: " # Добавлена", Path: path, IsNode: true})
		case ok && !yes && isObject1:
			diffs = append(diffs, Diff{DifTest: map[string]any{key: value1}, Message: " # Удалена", Path: path, IsNode: true})
		case ok && !yes && isSimple1:
			diffs = append(diffs, Diff{DifTest: map[string]any{key: value1}, Message: " # Удалена", Path: path})
		case !ok && yes && isSimple2:
			diffs = append(diffs, Diff{DifTest: map[string]any{key: value2}, Message: " # Добавлена", Path: path})
		case isSimple1 && isSimple2 && value1 != value2:
			diffs = append(diffs, Diff{DifTest: map[string]any{key: value1}, Message: " # Старое значение", Path: path})
			diffs = append(diffs, Diff{DifTest: map[string]any{key: value2}, Message: " # Новое значение", Path: path})
		case isSimple1 && isSimple2 && value1 == value2:
			diffs = append(diffs, Diff{DifTest: map[string]any{key: value1}, Message: " # Равны", Path: path})
		case isSimple1 && value2 == nil || isSimple2 && value1 == nil:
			diffs = append(diffs, Diff{DifTest: map[string]any{key: value1}, Message: " # Старое значение", Path: path})
			diffs = append(diffs, Diff{DifTest: map[string]any{key: value2}, Message: " # Новое значение", Path: path})
		default:
			diffs = append(diffs, Diff{DifTest: map[string]any{key: value1}, Message: " # Старое значение", Path: path})
			diffs = append(diffs, Diff{DifTest: map[string]any{key: value2}, Message: " # Новое значение", Path: path})
		}
	}
	return diffs
}

func IsObject(value interface{}) bool {
	_, ok := value.(map[string]interface{})
	return ok
}

func SimpleType(value interface{}) bool {
	switch value.(type) {
	case string, int, int8, int16, int32, int64, float32, float64, bool:
		return true
	default:
		return false
	}
}

// Рабочий код
//func Compare(obj1, obj2 map[string]interface{}, basePath string) []Diff {
//	var diffs []Diff
//
//	// Создаём множество ключей из обеих map
//	uniqueKeys := make(map[string]struct{}) // Убираю дубликаты ключей
//	for key := range obj1 {
//		uniqueKeys[key] = struct{}{}
//	}
//	for key := range obj2 {
//		uniqueKeys[key] = struct{}{}
//	}
//	commonKeys := make([]string, 0, len(uniqueKeys)) //Сразу аллоцирую память под известное количество элементов
//	for key := range uniqueKeys {
//		commonKeys = append(commonKeys, key)
//	}
//	// Сортируем ключи в алфавитном порядке
//	sort.Slice(commonKeys, func(i, j int) bool {
//		return commonKeys[i] < commonKeys[j]
//	})
//	repeatedKeys := make(map[string]bool) //для обработки повторяющихся ключей
//	for _, key := range commonKeys {
//		if repeatedKeys[key] {
//			continue
//		}
//		value1, ok := obj1[key]
//		value2, yes := obj2[key]
//
//		path := basePath + "." + key
//
//		if ok && yes && IsObject(value1) && IsObject(value2) {
//			node := make(map[string]any)
//
//			o1 := value1.(map[string]interface{})
//			o2 := value2.(map[string]interface{})
//			childDiff := Compare(o1, o2, path) // рекурсивный шаг
//			node[key] = childDiff
//			dif := Diff{
//				DifTest: map[string]any{
//					key: childDiff,
//				},
//				Path:    path,
//				IsNode:  true,
//				Message: " # Объекты",
//			}
//			diffs = append(diffs, dif)
//			repeatedKeys[key] = true
//			continue
//		}
//		isSimple1 := SimpleType(value1)
//		isSimple2 := SimpleType(value2)
//		isObject1 := IsObject(value1)
//		isObject2 := IsObject(value2)
//		switch {
//		case ok && yes && isSimple1 && isObject2:
//			diffs = append(diffs, Diff{DifTest: map[string]any{key: value1}, Message: " # Старое значение", Path: path})
//			diffs = append(diffs, Diff{DifTest: map[string]any{key: value2}, Message: " # Новое значение", Path: path, IsNode: true})
//		case ok && yes && isObject1 && isSimple2:
//			diffs = append(diffs, Diff{DifTest: map[string]any{key: value1}, Message: " # Старое значение", Path: path, IsNode: true})
//			diffs = append(diffs, Diff{DifTest: map[string]any{key: value2}, Message: " # Новое значение", Path: path})
//		case !ok && yes && isObject2:
//			diffs = append(diffs, Diff{DifTest: map[string]any{key: value2}, Message: " # Добавлена", Path: path, IsNode: true})
//		case ok && !yes && isObject1:
//			diffs = append(diffs, Diff{DifTest: map[string]any{key: value1}, Message: " # Удалена", Path: path, IsNode: true})
//		case ok && !yes && isSimple1:
//			diffs = append(diffs, Diff{DifTest: map[string]any{key: value1}, Message: " # Удалена", Path: path})
//		case !ok && yes && isSimple2:
//			diffs = append(diffs, Diff{DifTest: map[string]any{key: value2}, Message: " # Добавлена", Path: path})
//		case isSimple1 && isSimple2 && value1 != value2:
//			diffs = append(diffs, Diff{DifTest: map[string]any{key: value1}, Message: " # Старое значение", Path: path})
//			diffs = append(diffs, Diff{DifTest: map[string]any{key: value2}, Message: " # Новое значение", Path: path})
//		case isSimple1 && isSimple2 && value1 == value2:
//			diffs = append(diffs, Diff{DifTest: map[string]any{key: value1}, Message: " # Равны", Path: path})
//		case isSimple1 && value2 == nil || isSimple2 && value1 == nil:
//			diffs = append(diffs, Diff{DifTest: map[string]any{key: value1}, Message: " # Старое значение", Path: path})
//			diffs = append(diffs, Diff{DifTest: map[string]any{key: value2}, Message: " # Новое значение", Path: path})
//		default:
//			diffs = append(diffs, Diff{DifTest: map[string]any{key: value1}, Message: " # Старое значение", Path: path})
//			diffs = append(diffs, Diff{DifTest: map[string]any{key: value2}, Message: " # Новое значение", Path: path})
//		}
//	}
//	return diffs
//}
//
//func IsObject(value interface{}) bool {
//	_, ok := value.(map[string]interface{})
//	return ok
//}
//
//func SimpleType(value interface{}) bool {
//	switch value.(type) {
//	case string, int, int8, int16, int32, int64, float32, float64, bool:
//		return true
//	default:
//		return false
//	}
//}
