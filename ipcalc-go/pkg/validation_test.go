package pkg

import (
	"reflect"
	"testing"
)

func TestToBinary(t *testing.T) {
	cases := []struct {
		input  int
		output [8]int
	}{
		{input: 234, output: [8]int{0, 1, 0, 1, 0, 1, 1, 1}},
		{input: 45, output: [8]int{1, 0, 1, 1, 0, 1, 0, 0}},
		{input: 127, output: [8]int{1, 1, 1, 1, 1, 1, 1, 0}},
	}

	for _, tcase := range cases {
		bits := toBinary(tcase.input)
		if !reflect.DeepEqual(bits, tcase.output) {
			t.Fatalf("expected: %q, got: %q", tcase.output, bits)
		}
	}
}
func TestToDecimal(t *testing.T) {
	cases := []struct {
		output int
		input  [8]int
	}{
		{output: 15, input: [8]int{1, 1, 1, 1, 0, 0, 0, 0}},
		{output: 234, input: [8]int{0, 1, 0, 1, 0, 1, 1, 1}},
		{output: 45, input: [8]int{1, 0, 1, 1, 0, 1, 0, 0}},
		{output: 127, input: [8]int{1, 1, 1, 1, 1, 1, 1, 0}},
	}

	for _, tcase := range cases {
		num := toDecimal(tcase.input)
		if num != tcase.output {
			t.Fatalf("expected: %d, got: %q", tcase.output, num)
		}
	}
}
