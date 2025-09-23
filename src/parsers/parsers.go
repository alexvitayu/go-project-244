package parsers

import (
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

func ParseDataFromFiles(path1, path2 string) (res1, res2 map[string]interface{}, err error) {
	file1, err := os.ReadFile(path1)
	if err != nil {
		return nil, nil, errors.New("не удалось прочитать файл")
	}
	file2, err := os.ReadFile(path2)
	if err != nil {
		return nil, nil, errors.New("не удалось прочитать файл")
	}
	var data1 map[string]interface{}
	var data2 map[string]interface{}
	var ext1 = filepath.Ext(path1)
	var ext2 = filepath.Ext(path2)

	switch {
	case ext1 == ".json" && ext2 == ".json":
		if err := json.Unmarshal(file1, &data1); err != nil {
			return nil, nil, errors.New("не удалось преобразовать из json")
		}
		if err := json.Unmarshal(file2, &data2); err != nil {
			return nil, nil, errors.New("не удалось преобразовать из json")
		}
		return data1, data2, nil

	case ext1 == ".yaml" || ext1 == ".yml" && ext2 == ".yaml" || ext2 == ".yml":
		if err := yaml.Unmarshal(file1, &data1); err != nil {
			return nil, nil, errors.New("не удалось преобразовать из yaml")
		}
		if err := yaml.Unmarshal(file2, &data2); err != nil {
			return nil, nil, errors.New("не удалось преобразовать из yaml")
		}
		return data1, data2, nil

	case ext1 != ext2:
		return nil, nil, errors.New("разные форматы данных")

	default:
		return nil, nil, errors.New("неизвестный формат данных")
	}
}
