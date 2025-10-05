package code

import (
	"code/src/compare"
	"code/src/formatters"
	"code/src/parsers"
	"fmt"
	"path/filepath"
)

func GenDiff(path1, path2, format string) (string, error) {
	abs1, abs2, err := toAbsolutePath(path1, path2)
	if err != nil {
		return "", fmt.Errorf("toAbsolutePath: %w", err)
	}
	data1, data2, err := parsers.ParseDataFromFiles(abs1, abs2)
	if err != nil {
		return "", fmt.Errorf("parseDataFromFiles: %w", err)
	}
	basePath := ""
	diff := compare.Compare(data1, data2, basePath)
	formater := formatters.Format(format)
	str, err := formater.FormatDiff(diff)
	if err != nil {
		return "", fmt.Errorf("formatDiff: %w", err)
	}
	return str, nil
}

func toAbsolutePath(path1, path2 string) (string, string, error) {
	var abs1 string
	var abs2 string
	var err error
	if filepath.IsAbs(path1) && filepath.IsAbs(path2) {
		abs1 = path1
		abs2 = path2
		return abs1, abs2, nil
	}
	abs1, err = filepath.Abs(path1)
	if err != nil {
		return "", "", fmt.Errorf("не удалось преобразовать путь %s в абсолютный: %w", path1, err)
	}
	abs2, err = filepath.Abs(path2)
	if err != nil {
		return "", "", fmt.Errorf("не удалось преобразовать путь %s в абсолютный: %w", path2, err)
	}
	return abs1, abs2, nil
}
