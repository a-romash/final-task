package shuntingYard

import "fmt"

const (
	RPNTokenTypeOperand  = 1
	RPNTokenTypeOperator = 2
)

// RPNToken represents an abstract token object in RPN(Reverse Polish notation) which could either be an operator or operand.
type RPNToken struct {
	Type  int
	Value interface{}
}

// NewRPNOperandToken creates an instance of operand RPNToken with specified value.
func NewRPNOperandToken(val float64) *RPNToken {
	return NewRPNToken(val, RPNTokenTypeOperand)
}

// NewRPNOperatorToken creates an instance of operator RPNToken with specified value.
func NewRPNOperatorToken(val string) *RPNToken {
	return NewRPNToken(val, RPNTokenTypeOperator)
}

// NewRPNToken creates an instance of RPNToken with specified value and type.
func NewRPNToken(val interface{}, typ int) *RPNToken {
	return &RPNToken{Value: val, Type: typ}
}

// IsOperand determines whether a token is an operand with a specified value.
func (token *RPNToken) IsOperand(val float64) bool {
	return token.Type == RPNTokenTypeOperand && token.Value.(float64) == val
}

// IsOperator determines whether a token is an operator with a specified value.
func (token *RPNToken) IsOperator(val string) bool {
	return token.Type == RPNTokenTypeOperator && token.Value.(string) == val
}

// GetDescription returns a string that describes the token.
func (token *RPNToken) GetDescription() string {
	return fmt.Sprintf("(%d)%v", token.Type, token.Value)
}
