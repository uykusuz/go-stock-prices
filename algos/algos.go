package algos

import (
	"fmt"
	"go-stock-prices/model"
	"go-stock-prices/util"

	"github.com/pkg/errors"
)

func CalculatePerformance(
	candles []*model.Candle,
	symbol string,
	from *model.Date,
	to *model.Date,
) (*model.Performance, error) {
	simpleReturn, err := CalculateSimpleReturn(candles)
	if err != nil {
		return nil, errors.Wrap(err, "failed to compute simple return")
	}

	maxDrawdown := CalculateMaxDrawdown(candles)

	return &model.Performance{
		Symbol:       symbol,
		SimpleReturn: simpleReturn,
		MaxDrawdown:  maxDrawdown,
		From:         from,
		To:           to,
	}, nil
}

func CalculateSimpleReturn(candles []*model.Candle) (float64, error) {
	if len(candles) < 1 {
		return 0, nil
	}

	if util.EqualToEp(candles[0].Close, 0) {
		return 0, fmt.Errorf("cannot compute simple return: first price is zero")
	}

	simpleReturn := (candles[len(candles)-1].Close - candles[0].Close) / candles[0].Close
	return simpleReturn, nil
}

func CalculateMaxDrawdown(candles []*model.Candle) float64 {
	if len(candles) < 1 {
		return 0
	}

	maxDrawdown := float64(0)
	currentMax := candles[0].Close
	currentMin := currentMax

	for i, candle := range candles {
		if i == 0 {
			continue
		}

		if util.GreaterThanEp(candle.Close, currentMax) {
			currentMax = candle.Close
			currentMin = candle.Close
		} else if util.LessThanEp(candle.Close, currentMin) {
			currentMin = candle.Close
			currentDrawdown := (currentMax - currentMin) / currentMax

			if util.GreaterThanEp(currentDrawdown, maxDrawdown) {
				maxDrawdown = currentDrawdown
			}
		}
	}

	return maxDrawdown
}
