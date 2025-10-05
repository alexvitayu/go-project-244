package formatters

import (
	"code/src"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatDiffJson(t *testing.T) {
	path1 := filepath.Join("../../testdata/fixture/complicated/file1-1.yaml") //фикстура №1
	path2 := filepath.Join("../../testdata/fixture/complicated/file2-1.yaml") //фикстура №2
	path3 := filepath.Join("../../testdata/fixture/complicated/expect.json")  //фикстура №3

	diff, err := src.GenDiff(path1, path2)
	require.NoError(t, err) // если произошла ошибка в основном коде, то и нет смысла двигаться дальше
	formater := Format("json")
	jsonDiff, err := formater.FormatDiff(diff)
	require.NoError(t, err)
	expected, err := os.ReadFile(path3)
	require.NoError(t, err) // если произошла ошибка извлечения, то нет смысла идти дальше
	assert.Equal(t, string(expected), jsonDiff)
}
