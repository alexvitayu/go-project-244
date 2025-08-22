package code

import (
	"strings"
	"testing"
)

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

func TestGenDiff(t *testing.T) {
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
