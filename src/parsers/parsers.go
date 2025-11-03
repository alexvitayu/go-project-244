package parsers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func ParseDataFromFiles(path1, path2 string) (res1, res2 any, err error) {
	file1, err := os.ReadFile(path1)
	if err != nil {
		return nil, nil, fmt.Errorf("reading file %s: %w", path1, err)
	}
	file2, err := os.ReadFile(path2)
	if err != nil {
		return nil, nil, fmt.Errorf("reading file %s: %w", path2, err)
	}
	var data1 any
	var data2 any
	var ext1 = filepath.Ext(path1)
	var ext2 = filepath.Ext(path2)

	switch {
	case ext1 == ".json" && ext2 == ".json":
		if err := json.Unmarshal(file1, &data1); err != nil {
			return nil, nil, fmt.Errorf("json_unmarshal %v: %w", file1, err)
		}
		if err := json.Unmarshal(file2, &data2); err != nil {
			return nil, nil, fmt.Errorf("json_unmarshal %v: %w", file2, err)
		}
		return data1, data2, nil

	case ext1 == ".yaml" || ext1 == ".yml" && ext2 == ".yaml" || ext2 == ".yml":
		if err := yaml.Unmarshal(file1, &data1); err != nil {
			return nil, nil, fmt.Errorf("yaml_unmarshal %v: %w", file1, err)
		}
		if err := yaml.Unmarshal(file2, &data2); err != nil {
			return nil, nil, fmt.Errorf("yaml_unmarshal %v: %w", file2, err)
		}
		return data1, data2, nil

	case ext1 != ext2:
		return nil, nil, errors.New("mismatched file formats")

	default:
		return nil, nil, errors.New("unknown file formats")
	}
}
