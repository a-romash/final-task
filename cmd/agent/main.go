package main

import (
	"calculator/internal/agent"
	"calculator/pkg/config"
	rabbitmq "calculator/pkg/rabbitmq/agent"
)

func main() {
	config.Init()
	newAgent := agent.NewAgent(config.Config.AmountOfCalcs)
	rabbitmq.Init(newAgent)
}
