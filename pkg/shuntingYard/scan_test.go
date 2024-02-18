package shuntingYard

import (
	"reflect"
	"testing"
)

func TestScan(t *testing.T) {
	var expected, values []string
	var err error

	values, err = Scan("1+2-3")
	if err != nil {
		panic(err)
	}
	expected = []string{"1", "+", "2", "-", "3"}
	if !reflect.DeepEqual(values, expected) {
		t.Fatalf("Expected %v, Got %v", expected, values)
	}

	values, err = Scan("123^456")
	if err != nil {
		panic(err)
	}
	expected = []string{"123", "^", "456"}
	if !reflect.DeepEqual(values, expected) {
		t.Fatalf("Expected %v, Got %v", expected, values)
	}
}
