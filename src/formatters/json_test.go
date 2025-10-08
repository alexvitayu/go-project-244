package formatters_test

import (
	"code"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatDiffJson_Hexlet(t *testing.T) {
	path1 := filepath.Join("../../testdata/hexlet_testdata/file1.json")       //фикстура №1
	path2 := filepath.Join("../../testdata/hexlet_testdata/file2.json")       //фикстура №2
	path3 := filepath.Join("../../testdata/hexlet_testdata/result_json.json") //фикстура №3

	jsonDiff, err := code.GenDiff(path1, path2, "json")
	require.NoError(t, err) // если произошла ошибка в основном коде, то и нет смысла двигаться дальше
	expected, err := os.ReadFile(path3)
	require.NoError(t, err) // если произошла ошибка извлечения, то нет смысла идти дальше
	assert.Equal(t, string(expected), jsonDiff)
}
