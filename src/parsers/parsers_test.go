package parsers

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDataFromFiles(t *testing.T) {
	var testCases = []struct {
		name      string
		filename1 string
		filename2 string
		expected  [3]interface{}
	}{
		{
			name:      "BothFilesAreJson",
			filename1: filepath.Join("../../testdata/fixture/file1.json"),
			filename2: filepath.Join("../../testdata/fixture/file2.json"),
			expected: [3]interface{}{map[string]interface{}{"follow": false, "host": "hexlet.io", "proxy": "123.234.53.22", "timeout": float64(50)},
				map[string]interface{}{"host": "hexlet.io", "timeout": float64(20), "verbose": true}, nil},
		},
		{
			name:      "BothFilesAreYaml",
			filename1: filepath.Join("../../testdata/fixture/file1-1.yaml"),
			filename2: filepath.Join("../../testdata/fixture/file2-1.yaml"),
			expected: [3]interface{}{map[string]interface{}{"follow": false, "host": "hexlet.io", "proxy": "123.234.53.22", "timeout": 50},
				map[string]interface{}{"host": "hexlet.io", "timeout": 20, "verbose": true}, nil},
		},
		{
			name:      "WrongFormatFiles",
			filename1: filepath.Join("../../testdata/fixture/file1.pdf"),
			filename2: filepath.Join("../../testdata/fixture/file2.pdf"),
			expected:  [3]interface{}{nil, nil, errors.New("неизвестный формат данных")},
		},
		{
			name:      "DifferentFileFormats",
			filename1: filepath.Join("../../testdata/fixture/file1.json"),
			filename2: filepath.Join("../../testdata/fixture/file2.pdf"),
			expected:  [3]interface{}{nil, nil, errors.New("разные форматы данных")},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got1, got2, err := ParseDataFromFiles(tc.filename1, tc.filename2)
			fmt.Println(got1, got2)
			if got1 != nil && got2 != nil {
				assert.Equal(t, got1, tc.expected[0])
				assert.Equal(t, got2, tc.expected[1])
			} else {
				assert.Nil(t, got1)
				assert.Nil(t, got2)
			}
			expectedErr := tc.expected[2]
			if expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorAs(t, err, &expectedErr)
			}
		})
	}
}
