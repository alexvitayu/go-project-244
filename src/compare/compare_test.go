package compare_test

import (
	"code/src/compare"
	"code/src/parsers"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompare(t *testing.T) {
	var testCases = []struct {
		name  string
		path1 string
		path2 string
		path3 string
	}{
		{
			name:  "simpleJsonFiles",
			path1: "../../testdata/fixture/file1.json",
			path2: "../../testdata/fixture/file2.json",
			path3: "../../testdata/fixture/simpleDiff.txt",
		},
		{
			name:  "complicatedYamlFiles",
			path1: "../../testdata/fixture/complicated/file1-1.yaml",
			path2: "../../testdata/fixture/complicated/file2-1.yaml",
			path3: "../../testdata/fixture/complicated/complicatedDiff.txt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val1, val2, err := parsers.ParseDataFromFiles(tc.path1, tc.path2)
			require.NoError(t, err)
			got := compare.Compare(val1, val2, "")
			require.NoError(t, err)
			bytes, err := os.ReadFile(tc.path3)
			require.NoError(t, err)
			assert.Equal(t, string(bytes), fmt.Sprint(got))
		})
	}
}
