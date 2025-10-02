package compare

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type Diff struct {
	IsNode  bool
	Path    string
	Message string
	DifTest map[string]any
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
				Path:   path,
				IsNode: true,
			}
			diffs = append(diffs, dif)
			repeatedKeys[key] = true
			continue
		}
		isSimple1 := simpleType(value1)
		isSimple2 := simpleType(value2)
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
	SaveResult(diffs)
	return diffs
}

func IsObject(value interface{}) bool {
	_, ok := value.(map[string]interface{})
	return ok
}

func IsArray(value interface{}) bool {
	_, ok := value.([]interface{})
	return ok
}

func simpleType(value interface{}) bool {
	switch value.(type) {
	case string, int, int8, int16, int32, int64, float32, float64, bool:
		return true
	default:
		return false
	}
}

func SaveResult(diffs []Diff) {
	file, err := os.Create("compare.json")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	data, err := json.MarshalIndent(diffs, "", "")
	if err != nil {
		fmt.Println(err)
	}
	file.Write(data)
}

//func (df *Diff) PrintPath() string {
//	b := strings.Builder{}
//	for _, path := range df.Path {
//		b.WriteString(path + ".")
//	}
//	b.WriteString("/")
//	return b.String()
//}

//working variant
//type Diff struct {
//	//Path    []string
//	Message string
//	DifTest map[string]any
//}
//
//func Compare(obj1, obj2 map[string]interface{}) []Diff {
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
//
//		value1, ok := obj1[key]
//		value2, yes := obj2[key]
//
//		if IsObject(value2) && IsObject(value1) {
//			node := make(map[string]any)
//
//			o1 := value1.(map[string]interface{})
//			o2 := value2.(map[string]interface{})
//			childDiff := Compare(o1, o2) // рекурсивный шаг
//			node[key] = childDiff
//			dif := Diff{
//				DifTest: map[string]any{
//					key: childDiff,
//				},
//			}
//
//			diffs = append(diffs, dif)
//
//			repeatedKeys[key] = true
//			continue
//		}
//		switch {
//		case !ok && yes:
//			//strData2 := fmt.Sprintf(" + %v: %v", key, value2)
//			diffs = append(diffs, Diff{DifTest: map[string]any{key: value2}, Message: " # Добавлена"})
//		case ok && !yes:
//			//strData1 := fmt.Sprintf(" - %v: %v", key, value1)
//			diffs = append(diffs, Diff{DifTest: map[string]any{key: value1}, Message: " # Удалена"})
//
//		default:
//			isSimple1 := simpleType(value1)
//			isSimple2 := simpleType(value2)
//			// Оба простых типа и значения не равны
//			if isSimple1 && isSimple2 && value1 != value2 {
//				//strData1 := fmt.Sprintf(" - %v: %v", key, value1)
//				diffs = append(diffs, Diff{DifTest: map[string]any{key: value1}, Message: " # Старое значение"})
//				//strData2 := fmt.Sprintf(" + %v: %v", key, value2)
//				diffs = append(diffs, Diff{DifTest: map[string]any{key: value2}, Message: " # Новое значение"})
//				// Оба простых типа и значения равны
//			} else if isSimple1 && isSimple2 && value1 == value2 {
//				//strData1 := fmt.Sprintf("   %v: %v", key, value1)
//				diffs = append(diffs, Diff{DifTest: map[string]any{key: value1}, Message: ""})
//				// Если оба типа и не простые типы и не объекты
//			} else {
//				//strData1 := fmt.Sprintf(" - %v: %v", key, value1)
//				diffs = append(diffs, Diff{DifTest: map[string]any{key: value1}, Message: " # Старое значение"})
//				//strData2 := fmt.Sprintf(" + %v: %v", key, value2)
//				diffs = append(diffs, Diff{DifTest: map[string]any{key: value2}, Message: " # Новое значение"})
//			}
//		}
//		SaveResult(diffs)
//	}
//	return diffs
//}
//
//func IsObject(value interface{}) bool {
//	_, ok := value.(map[string]interface{})
//	return ok
//}
//
//func IsArray(value interface{}) bool {
//	_, ok := value.([]interface{})
//	return ok
//}
//
//func simpleType(value interface{}) bool {
//	switch value.(type) {
//	case string, int, int8, int16, int32, int64, float32, float64, bool:
//		return true
//	default:
//		return false
//	}
//}
//
//func SaveResult(diffs []Diff) {
//	file, err := os.Create("")
//	if err != nil {
//		fmt.Println(err)
//	}
//	defer file.Close()
//
//	data, err := json.MarshalIndent(diffs, "", "")
//	if err != nil {
//		fmt.Println(err)
//	}
//	file.Write(data)
//}
