package code

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenDiff(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name   string
		path1  string
		path2  string
		format string
		want   string
		err    bool
		errMsg string
	}{
		{"output_stylish",
			"testdata/fixture/complicated/file1-1.yaml",
			"testdata/fixture/complicated/file2-1.yaml",
			"stylish",
			"testdata/fixture/complicated/expectStylish.txt",
			false,
			"",
		},
		{"output_plain",
			"testdata/fixture/complicated/file1-1.yaml",
			"testdata/fixture/complicated/file2-1.yaml",
			"plain",
			"testdata/fixture/complicated/expectPlain.txt",
			false,
			"",
		},
		{"output_json",
			"testdata/hexlet_testdata/file1.json",
			"testdata/hexlet_testdata/file2.json",
			"json",
			"testdata/hexlet_testdata/result_json.json",
			false,
			"",
		},
		{"error_incorrect_path",
			"",
			"testdata/hexlet_testdata/file2.json",
			"json",
			"",
			true,
			"toAbsolutePath: некорректный ввод путей",
		},
		{"error_different_file_formats",
			"testdata/hexlet_testdata/file1.yml",
			"testdata/hexlet_testdata/file2.json",
			"json",
			"",
			true,
			"parseDataFromFiles: разные форматы данных",
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := GenDiff(tc.path1, tc.path2, tc.format)
			if tc.err {
				require.EqualError(t, err, tc.errMsg)
			} else {
				require.NoError(t, err)
				bytes, err := os.ReadFile(tc.want)
				require.NoError(t, err)
				assert.Equal(t, got, string(bytes))
			}
		})
	}
}

func TestToAbsolutePath(t *testing.T) {
	t.Parallel()
	temp := t.TempDir()
	path1 := filepath.Join(temp, "file1.json")
	path2 := filepath.Join(temp, "file2.json")
	want1, err := filepath.Abs(path1)
	require.NoError(t, err)
	want2, err := filepath.Abs(path2)
	require.NoError(t, err)
	got1, got2, err := toAbsolutePath(path1, path2)
	require.NoError(t, err)
	assert.Equal(t, got1, want1)
	assert.Equal(t, got2, want2)
}
