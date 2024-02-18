package shuntingYard

import "testing"

func TestEvaluate(t *testing.T) {
	var tokens []*RPNToken
	var expected, got float64
	var err error

	tokens = []*RPNToken{NewRPNOperandToken(1), NewRPNOperandToken(32), NewRPNOperandToken(3), NewRPNOperatorToken("-"), NewRPNOperatorToken("*")}
	expected = 29
	got, err = Evaluate(tokens)
	if err != nil {
		panic(err)
	}
	if got != expected {
		t.Fatalf("Expected %v, Got %v.", expected, got)
	}

	tokens = []*RPNToken{NewRPNOperandToken(1), NewRPNOperandToken(2), NewRPNOperandToken(2), NewRPNOperatorToken("+"), NewRPNOperandToken(2), NewRPNOperandToken(4), NewRPNOperatorToken("*"), NewRPNOperatorToken("/"), NewRPNOperatorToken("-")}
	if err != nil {
		panic(err)
	}

	expected = 0.5
	got, err = Evaluate(tokens)
	if err != nil {
		panic(err)
	}
	if got != expected {
		t.Fatalf("Expected %v, Got %v.", expected, got)
	}
}
