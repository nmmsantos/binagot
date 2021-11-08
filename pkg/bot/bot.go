package bot

import (
	"github.com/nmmsantos/binagot/pkg/account"
	"github.com/nmmsantos/binagot/pkg/log"
)

const (
	Buy  Operation = "BUY"
	Sell           = "SELL"
)

type (
	Bot interface {
		GetAvailableSymbols() ([]string, error)
		GetPrices() (BrokerPrices, error)
		ShouldSave() bool
	}

	BrokerPrices map[string]float64

	Trader interface {
		Trade(*account.Accout) []OperationRequest
	}

	OperationRequest struct {
		Operation Operation
		Symbol    *account.Symbol
		Buckets   int
	}

	Operation string
)

var logger = log.GetLogger()

func Run(bot Bot, trader Trader, configFile string) error {
	logger.Info("starting bot")

	acc := &account.Accout{}

	logger.Info("loading account")
	if err := acc.Load(configFile); err != nil {
		return err
	}

	availableSymbolsMap := map[string]string{}

	logger.Info("fetching available symbols")
	if availableSymbols, err := bot.GetAvailableSymbols(); err != nil {
		return err
	} else {
		for _, s := range availableSymbols {
			availableSymbolsMap[s] = s
		}
	}

	loadedSymbolsMap := map[string]*account.Symbol{}

	for _, s := range acc.Symbols {
		loadedSymbolsMap[s.Name] = s
	}

	acc.Symbols = []*account.Symbol{}

	logger.Info("generating possible symbols")
	for _, base := range acc.Assets {
		for _, quote := range acc.Assets {
			s := base.Name + quote.Name

			if _, exists := availableSymbolsMap[s]; exists {
				logger.Info("trading symbol ", s)

				ts := &account.Symbol{
					Name:       s,
					BaseAsset:  base.Name,
					QuoteAsset: quote.Name,
				}

				if ls, exists := loadedSymbolsMap[s]; exists {
					ts.Price, ts.RefPrice, ts.RefDate = ls.Price, ls.RefPrice, ls.RefDate
				}

				acc.Symbols = append(acc.Symbols, ts)
			}
		}
	}

	if bot.ShouldSave() {
		logger.Info("saving account")

		if err := acc.Save(configFile); err != nil {
			logger.Error("couldn't save file", configFile)
		}
	}

	initialAccountValue, fiat := acc.GetValue()
	transactionFee := fiat.Ammount / float64(fiat.Buckets) * 0.0075
	transactionCount := 0

	defer func() {
		value, _ := acc.GetValue()
		logger.Infof("initial account value: %f %s", initialAccountValue, fiat.Name)
		logger.Infof("transactions: %d * %f %s fee", transactionCount, transactionFee, fiat.Name)
		logger.Infof("account value: %f %s", value-float64(transactionCount)*transactionFee, fiat.Name)
	}()

	for {
		if prices, err := bot.GetPrices(); err != nil {
			logger.Error(err)
		} else if prices == nil {
			return nil
		} else {
			for _, s := range acc.Symbols {
				if p, exists := prices[s.Name]; exists {
					s.Price = p
				}
			}

			reqs := trader.Trade(acc)

			for _, r := range reqs {
				logger.Infof("request to %s %d bucket(s) of %s at %f %s", r.Operation, r.Buckets, r.Symbol.BaseAsset, r.Symbol.Price, r.Symbol.QuoteAsset)

				var base, quote *account.Asset

				for _, a := range acc.Assets {
					if a.Name == r.Symbol.BaseAsset {
						base = a
					}
					if a.Name == r.Symbol.QuoteAsset {
						quote = a
					}
				}

				if r.Operation == Buy {
					if r.Buckets > quote.Buckets {
						logger.Infof("not enough bucket(s) of %s", quote.Name)
						continue
					} else {
						quoteAmmount := quote.Ammount / float64(quote.Buckets) * float64(r.Buckets)
						baseAmmount := quoteAmmount / r.Symbol.Price
						quote.Ammount -= quoteAmmount
						quote.Buckets -= r.Buckets
						base.Ammount += baseAmmount
						base.Buckets += r.Buckets
						logger.Infof("bought %f %s with %f %s", baseAmmount, base.Name, quoteAmmount, quote.Name)
					}
				} else if r.Operation == Sell {
					if r.Buckets > base.Buckets {
						logger.Infof("not enough bucket(s) of %s", base.Name)
						continue
					} else {
						baseAmmount := base.Ammount / float64(base.Buckets) * float64(r.Buckets)
						quoteAmmount := baseAmmount * r.Symbol.Price
						base.Ammount -= baseAmmount
						base.Buckets -= r.Buckets
						quote.Ammount += quoteAmmount
						quote.Buckets += r.Buckets
						logger.Infof("sold %f %s at %f %s", baseAmmount, base.Name, quoteAmmount, quote.Name)
					}
				}

				transactionCount++
				logger.Infof("account summary: %s", acc.GetSummary())
			}
		}
	}
}
