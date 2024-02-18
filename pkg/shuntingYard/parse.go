package shuntingYard

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// precedence of operators
var priorities map[string]float64

// associativities of operators
var associativities map[string]bool

func init() {
	priorities = make(map[string]float64, 0)
	associativities = make(map[string]bool, 0)

	priorities["+"] = 0
	priorities["-"] = 0
	priorities["*"] = 1
	priorities["/"] = 1
	priorities["^"] = 2

	// if not set, associativity will be false(left-associative)
}

// Parse parses an array of token strings and returns an array of abstract tokens
// using Shunting-yard algorithm.
func Parse(tokens []string) ([]*RPNToken, error) {
	var ret []*RPNToken

	var operators []string
	for _, token := range tokens {
		operandToken := tryGetOperand(token)
		if operandToken != nil {
			ret = append(ret, operandToken)
		} else {
			// check parentheses
			if token == "(" {
				operators = append(operators, token)
			} else if token == ")" {
				foundLeftParenthesis := false
				// pop until "(" is fouund
				for len(operators) > 0 {
					oper := operators[len(operators)-1]
					operators = operators[:len(operators)-1]

					if oper == "(" {
						foundLeftParenthesis = true
						break
					} else {
						ret = append(ret, NewRPNOperatorToken(oper))
					}
				}
				if !foundLeftParenthesis {
					return nil, errors.New("Mismatched parentheses found")
				}
			} else {
				// operator priority and associativity
				priority, ok := priorities[token]
				if !ok {
					return nil, fmt.Errorf("Unknown operator: %v", token)
				}
				rightAssociative := associativities[token]

				for len(operators) > 0 {
					top := operators[len(operators)-1]

					if top == "(" {
						break
					}

					prevPriority := priorities[top]

					if (rightAssociative && priority < prevPriority) || (!rightAssociative && priority <= prevPriority) {
						// pop current operator
						operators = operators[:len(operators)-1]
						ret = append(ret, NewRPNOperatorToken(top))
					} else {
						break
					}
				} // end of for len(operators) > 0

				operators = append(operators, token)
			} // end of if token == "("
		} // end of if isOperand(token)
	} // end of for _, token := range tokens

	// process remaining operators
	for len(operators) > 0 {
		// pop
		operator := operators[len(operators)-1]
		operators = operators[:len(operators)-1]

		if operator == "(" {
			return nil, errors.New("Mismatched parentheses found")
		}
		ret = append(ret, NewRPNOperatorToken(operator))
	}
	return ret, nil
}

// tryGetOperand determines whether a given string is an operand, if it is, an RPN operand token will be returned, otherwise nil.
func tryGetOperand(str string) *RPNToken {
	value, err := strconv.ParseFloat(strings.TrimSpace(str), 64)
	if err != nil {
		return nil
	}
	return NewRPNOperandToken(value)
}
