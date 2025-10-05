package formatters_test

import (
	"code"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Сравнение двух .yaml файлов
func TestFormatDiffStylish(t *testing.T) {
	path1 := filepath.Join("../../testdata/fixture/complicated/file1-1.yaml")      //фикстура №1
	path2 := filepath.Join("../../testdata/fixture/complicated/file2-1.yaml")      //фикстура №2
	path3 := filepath.Join("../../testdata/fixture/complicated/expectStylish.txt") //фикстура №3
	str, err := code.GenDiff(path1, path2, "stylish")
	require.NoError(t, err)
	bytes, err := os.ReadFile(path3)
	require.NoError(t, err)             // если произошла ошибка чтения, то нет смысла идти дальше
	assert.Equal(t, string(bytes), str) // утверждаем, что ожидаемая строка равна полученной строке
}

// Сравнение плоских структур
func TestFormatDiff(t *testing.T) {
	path1 := filepath.Join("../../testdata/fixture/file1.yml")  //фикстура №1
	path2 := filepath.Join("../../testdata/fixture/file2.yml")  //фикстура №2
	path3 := filepath.Join("../../testdata/fixture/result.txt") //фикстура №3

	str, err := code.GenDiff(path1, path2, "stylish")
	require.NoError(t, err)
	bytes, err := os.ReadFile(path3)
	require.NoError(t, err)             // если произошла ошибка чтения, то нет смысла идти дальше
	assert.Equal(t, string(bytes), str) // утверждаем, что ожидаемая строка равна полученной строке
}
