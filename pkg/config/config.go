package config

import (
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type IConfig struct {
	ServerAddr       string
	AgentTimeout     int
	AgentPing        int
	AgentResolveTime map[string]int
	RabbitMQURL      string
	AmountOfCalcs    int
}

var Config *IConfig
var once sync.Once

func Init() error {
	var err error = nil
	once.Do(func() {
		err = godotenv.Load()

		Config = &IConfig{
			ServerAddr:  os.Getenv("SERVER_ADDR"),
			RabbitMQURL: os.Getenv("RABBITMQ_URL"),
		}

		Config.AgentTimeout, err = strconv.Atoi(os.Getenv("AGENT_TIMEOUT"))

		Config.AgentPing, err = strconv.Atoi(os.Getenv("AGENT_PING"))

		var plusResolveTime, minusResolveTime, multiplyResolveTime, divideResolveTime int

		plusResolveTime, err = strconv.Atoi(os.Getenv("PLUS_RESOLVE_TIME"))
		minusResolveTime, err = strconv.Atoi(os.Getenv("MINUS_RESOLVE_TIME"))
		multiplyResolveTime, err = strconv.Atoi(os.Getenv("MULTIPLY_RESOLVE_TIME"))
		divideResolveTime, err = strconv.Atoi(os.Getenv("DIVIDE_RESOLVE_TIME"))

		Config.AgentResolveTime = map[string]int{
			"+": plusResolveTime,
			"-": minusResolveTime,
			"*": multiplyResolveTime,
			"/": divideResolveTime,
		}

		Config.AmountOfCalcs, err = strconv.Atoi(os.Getenv("AMOUNT_OF_CALCS"))
	})
	return err
}
