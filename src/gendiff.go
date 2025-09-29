package src

import (
	"code/src/compare"
	"code/src/formatter"
	"code/src/parsers"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

func GenDiff(path1, path2, format string) (string, error) {
	abs1, abs2 := toAbsolutePath(path1, path2)
	data1, data2, err := parsers.ParseDataFromFiles(abs1, abs2)
	if err != nil {
		return "", fmt.Errorf("parseDataFromFiles: %w", err)
	}
	diff := compare.Compare(data1, data2)

	str := formatter.FormatDiff(diff)

	return str, nil
}

func toAbsolutePath(path1, path2 string) (p1, p2 string) {
	var abs1 string
	var abs2 string
	if filepath.IsAbs(path1) && filepath.IsAbs(path2) {
		abs1 = path1
		abs2 = path2
		return abs1, abs2
	}
	abs1, _ = filepath.Abs(path1)
	abs2, _ = filepath.Abs(path2)

	return abs1, abs2
}

func genDiff(data1, data2 map[string]interface{}) string {

	// Создаём множество ключей из обеих map
	commonKeys := make([]string, 0, (len(data1) + len(data2))) //Сразу аллоцирую память под известное количество элементов
	for key := range data1 {
		commonKeys = append(commonKeys, key)
	}
	for key := range data2 {
		commonKeys = append(commonKeys, key)
	}

	// Сортируем ключи а алфавитном порядке
	sort.Slice(commonKeys, func(i, j int) bool {
		return commonKeys[i] < commonKeys[j]
	})
	// Результирующий массив
	resSlice := make([]string, 0)

	// Результирующая строка
	resStr := strings.Builder{}

	// Итеративно пройдёмся по commonKeys
	repeatedKeys := make(map[string]bool) //для обработки повторяющихся ключей
	for _, key := range commonKeys {
		var strData1 string
		var strData2 string
		value1, ok := data1[key]
		value2, yes := data2[key]
		switch {
		case ok && yes && value1 != value2:
			_, ok := repeatedKeys[key]
			if !ok {
				strData1 = fmt.Sprintf(" - %v: %v", key, value1)
				resSlice = append(resSlice, strData1)
				//resStr.WriteString(fmt.Sprintf(" - %v: %v\n", key, value1))
				repeatedKeys[key] = true
			} else {
				strData2 = fmt.Sprintf(" + %v: %v", key, value2)
				resSlice = append(resSlice, strData2)
				//resStr.WriteString(fmt.Sprintf(" + %v: %v\n", key, value2))
			}
		case ok && value2 != value1:
			strData1 = fmt.Sprintf(" - %v: %v", key, value1)
			resSlice = append(resSlice, strData1)
			//resStr.WriteString(fmt.Sprintf(" - %v: %v\n", key, value1))
		case yes && value2 != value1:
			strData2 = fmt.Sprintf(" + %v: %v", key, value2)
			resSlice = append(resSlice, strData2)
			//resStr.WriteString(fmt.Sprintf(" + %v: %v\n", key, value2))
		case ok && yes && value1 == value2:
			_, ok := repeatedKeys[key]
			if !ok {
				strData1 = fmt.Sprintf("   %v: %v", key, value1)
				resSlice = append(resSlice, strData1)
				//resStr.WriteString(fmt.Sprintf("   %v: %v\n", key, value1))
				repeatedKeys[key] = true
			} else {
				continue
			}
		}
	}
	for i, v := range resSlice {
		resStr.WriteString(v)
		if i < len(resSlice)-1 {
			resStr.WriteString("\n")
		}
	}
	return resStr.String()
}
