package shuntingYard

import "testing"

func TestParse(t *testing.T) {
	tokens, err := Parse([]string{"1", "-", "32", "*", "3"})
	if err != nil {
		panic(err)
	}
	if !tokens[0].IsOperand(1) {
		t.Fatalf("Expected '1', Got %v", tokens[0].GetDescription())
	}
	if !tokens[1].IsOperand(32) {
		t.Fatalf("Expected '32', Got %v", tokens[1].GetDescription())
	}
	if !tokens[2].IsOperand(3) {
		t.Fatalf("Expected '3', Got %v", tokens[2].GetDescription())
	}
	if !tokens[3].IsOperator("*") {
		t.Fatalf("Expected '*', Got %v", tokens[3].GetDescription())
	}
	if !tokens[4].IsOperator("-") {
		t.Fatalf("Expected '-', Got %v", tokens[4].GetDescription())
	}
}
