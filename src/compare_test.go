package src

import (
	"code/src/parsers"
	"fmt"
	"path/filepath"
	"testing"
)

func TestCompare(t *testing.T) {
	path1 := filepath.Join("../testdata/fixture/complicated/file1-1.yaml") //фикстура №1
	path2 := filepath.Join("../testdata/fixture/complicated/file2-1.yaml") //фикстура №2
	//path1 := filepath.Join("../testdata/fixture/file1.yml")
	//path2 := filepath.Join("../testdata/fixture/file2.yml")
	obj1, obj2, _ := parsers.ParseDataFromFiles(path1, path2)

	res := Compare(obj1, obj2, "")
	fmt.Println(res)
}
