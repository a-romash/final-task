package agent

import (
	"calculator/internal/model"
	"log"
	"time"

	"calculator/pkg/shuntingYard"
)

type Calculator struct {
	taskChan   chan *model.ExpressionPart
	expression *model.ExpressionPart
	isBusy     bool
	id         int
}

func NewCalculator(i int) *Calculator {
	c := &Calculator{
		taskChan: make(chan *model.ExpressionPart),
		isBusy:   false,
		id:       i,
	}

	c.Start()

	return c
}

func (c *Calculator) AddTask(task *model.ExpressionPart) bool {
	if c.isBusy {
		return false
	}

	c.taskChan <- task
	return true
}

func (c *Calculator) Start() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
				c.Start()
			}
		}()

		for {
			task := <-c.taskChan
			c.expression = task
			c.isBusy = true

			c.SolveExpression(task)

			c.expression = nil
			c.isBusy = false
		}
	}()
}

func (c *Calculator) SolveExpression(expr *model.ExpressionPart) {
	time.Sleep(time.Duration(expr.Duration) * time.Second)

	if result, err := shuntingYard.Evaluate([]*shuntingYard.RPNToken{expr.FirstOperand, expr.SecondOperand, expr.Operation}); err == nil {
		tokenizedResult := shuntingYard.NewRPNOperandToken(result)
		expr.Result <- tokenizedResult
	}
	expr.Result <- shuntingYard.NewRPNOperandToken(0)
}
