package parsers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDataFromFiles(t *testing.T) {
	var testCases = []struct {
		name      string
		path1     string
		path2     string
		expected1 any
		expected2 any
		isErr     bool
		errText   string
	}{
		{
			name:      "yaml_files",
			path1:     "../../testdata/fixture/complicated/file1-1.yaml",
			path2:     "../../testdata/fixture/complicated/file2-1.yaml",
			expected1: `map[common:map[setting1:Value 1 setting2:200 setting3:true setting6:map[doge:map[wow:] key:value]] group1:map[baz:bas foo:bar nest:map[key:value]] group2:map[abc:12345 deep:map[id:45]]]`,
			expected2: `map[common:map[follow:false setting1:Value 1 setting3:<nil> setting4:blah blah setting5:map[key5:value5] setting6:map[doge:map[wow:so much] key:value ops:vops]] group1:map[baz:bars foo:bar nest:str] group3:map[deep:map[id:map[number:45]] fee:100500]]`,
			isErr:     false,
		},
		{
			name:      "yml_files",
			path1:     "../../testdata/fixture/file1.yml",
			path2:     "../../testdata/fixture/file2.yml",
			expected1: `map[follow:false host:hexlet.io proxy:123.234.53.22 timeout:50]`,
			expected2: `map[host:hexlet.io timeout:20 verbose:true]`,
			isErr:     false,
		},
		{
			name:      "json_files",
			path1:     "../../testdata/fixture/file1.json",
			path2:     "../../testdata/fixture/file2.json",
			expected1: `map[follow:false host:hexlet.io proxy:123.234.53.22 timeout:50]`,
			expected2: `map[host:hexlet.io timeout:20 verbose:true]`,
			isErr:     false,
		},
		{
			name:      "different_file_formats",
			path1:     "../../testdata/fixture/file1.yml",
			path2:     "../../testdata/fixture/file2.json",
			expected1: nil,
			expected2: nil,
			isErr:     true,
			errText:   "разные форматы данных",
		},
		{
			name:      "unknown_file_formats",
			path1:     "../../testdata/fixture/file1.pdf",
			path2:     "../../testdata/fixture/file2.pdf",
			expected1: nil,
			expected2: nil,
			isErr:     true,
			errText:   "неизвестный формат данных",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got1, got2, err := ParseDataFromFiles(tc.path1, tc.path2)
			if tc.isErr {
				require.Error(t, err)
				assert.Equal(t, tc.errText, err.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected1, fmt.Sprint(got1))
				assert.Equal(t, tc.expected2, fmt.Sprint(got2))
			}
		})
	}
}
