package main

import (
	"github.com/nmmsantos/binagot/pkg/bot"
	"github.com/nmmsantos/binagot/pkg/log"
	"github.com/nmmsantos/binagot/pkg/simulation"
	"github.com/nmmsantos/binagot/pkg/trading"
)

var logger = log.GetLogger()

func main() {
	if b, err := simulation.NewSimulationBot("simulation.csv"); err != nil {
		logger.Fatal(err)
	} else {
		if err := bot.Run(b, trading.NewPositionSizingTrader(0.1, 0.1), "config.yaml"); err != nil {
			logger.Fatal(err)
		}
	}
}
