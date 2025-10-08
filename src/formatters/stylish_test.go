package formatters_test

import (
	"code"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Сравнение двух .yaml файлов
func TestFormatDiffStylish(t *testing.T) {
	var testCases = []struct {
		name  string
		path1 string
		path2 string
		path3 string
	}{
		{
			name:  "complicatedYamlFiles",
			path1: "../../testdata/fixture/complicated/file1-1.yaml",
			path2: "../../testdata/fixture/complicated/file2-1.yaml",
			path3: "../../testdata/fixture/complicated/expectStylish.txt",
		},
		{
			name:  "flatYamlFiles",
			path1: "../../testdata/fixture/file1.yml",
			path2: "../../testdata/fixture/file2.yml",
			path3: "../../testdata/fixture/result.txt",
		},
		{
			name:  "hexletJsonFiles",
			path1: "../../testdata/hexlet_testdata/file1.json",
			path2: "../../testdata/hexlet_testdata/file2.json",
			path3: "../../testdata/hexlet_testdata/result_stylish.txt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := code.GenDiff(tc.path1, tc.path2, "stylish")
			require.NoError(t, err)
			bytes, err := os.ReadFile(tc.path3)
			require.NoError(t, err)
			assert.Equal(t, string(bytes), got)
		})
	}
}
