package account

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type (
	Accout struct {
		Assets  []*Asset  `validate:"required,dive"`
		Symbols []*Symbol `validate:"dive"`
	}

	Asset struct {
		Name    string  `validate:"uppercase"`
		Buckets int     `validate:"gte=0"`
		Ammount float64 `validate:"gte=0"`
	}

	Symbol struct {
		Name       string  `validate:"uppercase"`
		BaseAsset  string  `validate:"uppercase"`
		QuoteAsset string  `validate:"uppercase"`
		Price      float64 `validate:"gte=0"`
		RefPrice   float64 `validate:"gte=0"`
		RefDate    time.Time
	}
)

var validate *validator.Validate = validator.New()

func (a *Accout) Load(filename string) error {
	if file, err := os.OpenFile(filename, os.O_RDONLY, 0644); err != nil {
		return err
	} else {
		defer file.Close()

		dec := yaml.NewDecoder(file)

		if err := dec.Decode(a); err != nil {
			return err
		}
	}

	if err := validate.Struct(a); err != nil {
		return err
	}

	return nil
}

func (a *Accout) Save(filename string) error {
	if file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
		return err
	} else {
		defer file.Close()

		enc := yaml.NewEncoder(file)
		defer enc.Close()

		enc.SetIndent(2)

		if err := enc.Encode(a); err != nil {
			return err
		}
	}

	return nil
}

func (a *Accout) GetSummary() string {
	sb := strings.Builder{}

	for _, a := range a.Assets {
		fmt.Fprintf(&sb, "%s: %f/%d  ", a.Name, a.Ammount, a.Buckets)
	}

	return strings.TrimSpace(sb.String())
}

func (a *Accout) GetValue() (float64, *Asset) {
	if len(a.Assets) <= 0 {
		return 0, nil
	}

	fiat := a.Assets[0]
	value := 0.0

	for _, as := range a.Assets {
		if as.Name == fiat.Name {
			value += as.Ammount
		} else {
			for _, s := range a.Symbols {
				if s.QuoteAsset == fiat.Name && s.BaseAsset == as.Name {
					value += as.Ammount * s.Price
				}
			}
		}
	}

	return value, fiat
}
