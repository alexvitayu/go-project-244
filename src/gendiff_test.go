package src

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// мой тест с использованием фикстур
func TestGenDiffWithFixturesJson(t *testing.T) {
	path1 := filepath.Join("../testdata/fixture/file1.json") //фикстура №1
	path2 := filepath.Join("../testdata/fixture/file2.json") //фикстура №2
	path3 := filepath.Join("../testdata/fixture/result.txt") //фикстура №3

	str, err := GenDiff(path1, path2, "stylish")
	require.NoError(t, err) // если произошла ошибка в основном коде, то и нет смысла двигаться дальше

	bytes, err := os.ReadFile(path3)
	require.NoError(t, err)             // если произошла ошибка чтения, то нет смысла идти дальше
	assert.Equal(t, string(bytes), str) // утверждаем, что ожидаемая строка равна полученной строке
}

func TestGenDiffWithFixturesYaml(t *testing.T) {
	path1 := filepath.Join("../testdata/fixture/file1.yml")  //фикстура №1
	path2 := filepath.Join("../testdata/fixture/file2.yml")  //фикстура №2
	path3 := filepath.Join("../testdata/fixture/result.txt") //фикстура №3

	str, err := GenDiff(path1, path2, "stylish")
	require.NoError(t, err) // если произошла ошибка в основном коде, то и нет смысла двигаться дальше

	bytes, err := os.ReadFile(path3)
	require.NoError(t, err)             // если произошла ошибка чтения, то нет смысла идти дальше
	assert.Equal(t, string(bytes), str) // утверждаем, что ожидаемая строка равна полученной строке
}

// мой тест до использовния фикстур
func TestGenDiff(t *testing.T) {
	var testCases = []struct {
		name   string
		data1  map[string]interface{}
		data2  map[string]interface{}
		expect string
	}{
		{name: "test1",
			data1: map[string]interface{}{
				"host":    "hexlet.io",
				"timeout": 50,
				"proxy":   "123.234.53.22",
				"follow":  false,
			},
			data2: map[string]interface{}{
				"timeout": 20,
				"verbose": true,
				"host":    "hexlet.io",
			},
			expect: " - follow: false\n   host: hexlet.io\n - proxy: 123.234.53.22\n - timeout: 50\n + timeout: 20\n + verbose: true",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := genDiff(tc.data1, tc.data2)
			res := strings.TrimSpace(got)
			want := strings.TrimSpace(tc.expect)
			if res != want {
				t.Errorf("ожидалось %s, получили %s", tc.expect, got)
			}
		})
	}
}
