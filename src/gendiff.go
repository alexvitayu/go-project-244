package code

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func GenDiff(path1, path2, format string) (string, error) {
	abs1, abs2 := Path(path1, path2)
	ParseDataFromFiles(abs1, abs2)
	return "", nil
}

func Path(path1, path2 string) (p1, p2 string) {
	var abs1 string
	var abs2 string
	if filepath.IsAbs(path1) && filepath.IsAbs(path2) {
		abs1 = path1
		abs2 = path2
		return abs1, abs2
	}
	abs1, _ = filepath.Abs(path1)
	abs2, _ = filepath.Abs(path2)

	return abs1, abs2
}

func ParseDataFromFiles(abs1, abs2 string) {
	if filepath.Ext(abs1) == ".json" && filepath.Ext(abs2) == ".json" {
		data1 := orderedmap.New[string, interface{}]()
		data2 := orderedmap.New[string, interface{}]()
		dataFile1, _ := os.ReadFile(abs1)
		dataFile2, _ := os.ReadFile(abs2)
		json.Unmarshal(dataFile1, &data1)
		json.Unmarshal(dataFile2, &data2)
		for pair := data1.Oldest(); pair != nil; pair = pair.Next() {
			fmt.Printf("%v : %v\n", pair.Key, pair.Value)
		}
		for pair := data2.Oldest(); pair != nil; pair = pair.Next() {
			fmt.Printf("%v : %v\n", pair.Key, pair.Value)
		}
	}

}
