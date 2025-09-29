package formatter

import (
	"code/src/compare"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"
)

func TestFormatDiff(t *testing.T) {
	file, err := os.ReadFile("../../compare.json")
	if err != nil {
		t.Fatal(errors.New("не удалось прочитать файл"))
	}
	var data []compare.Diff
	err = json.Unmarshal(file, &data)

	str := FormatDiff(data)
	fmt.Println(str)
}
