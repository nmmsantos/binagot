package simulation

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/nmmsantos/binagot/pkg/bot"
)

type SimulationBot struct {
	symbols []string
	prices  [][]string
	current int
}

func NewSimulationBot(simulationFile string) (*SimulationBot, error) {
	b := &SimulationBot{}

	if file, err := os.OpenFile(simulationFile, os.O_RDONLY, 0644); err != nil {
		return b, err
	} else {
		defer file.Close()

		if data, err := csv.NewReader(file).ReadAll(); err != nil {
			return b, nil
		} else {
			b.symbols = data[0]
			b.prices = data[1:]
		}
	}

	return b, nil
}

func (b *SimulationBot) GetAvailableSymbols() ([]string, error) {
	return b.symbols, nil
}

func (b *SimulationBot) GetPrices() (bot.BrokerPrices, error) {
	if b.current >= len(b.prices) {
		return nil, nil
	}

	defer func() {
		b.current++
	}()

	prices := make(bot.BrokerPrices, len(b.symbols))

	for i, s := range b.symbols {
		if p, err := strconv.ParseFloat(b.prices[b.current][i], 64); err != nil {
			return prices, err
		} else {
			prices[s] = p
		}
	}

	return prices, nil
}

func (b *SimulationBot) ShouldSave() bool {
	return false
}
