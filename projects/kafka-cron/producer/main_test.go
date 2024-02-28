package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestParseCronFile(t *testing.T) {
	var tests = []struct {
		input    string
		expected []Job
	}{
		{
			"0 0 1 1 0 /bin/test",
			[]Job{
				{Command: "/bin/test", CronSpec: CronSpec{1, 1, 2, 2, 1}},
			},
		},
		{
			"0 5-10 * 1 0 /bin/test",
			[]Job{
				{Command: "/bin/test", CronSpec: CronSpec{1, 2016, 9223372041149743102, 2, 1}},
			},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			ans, err := parseCronFile(strings.NewReader(tt.input))
			if !reflect.DeepEqual(ans, tt.expected) || err != nil {
				t.Errorf("got %v, expected %v", ans, tt.expected)
			}
		})
	}
}
