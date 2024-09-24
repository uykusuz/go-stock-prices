package ports

import (
	"context"
	"go-stock-prices/model"
)

type FinancialDataRepository interface {
	GetCandles(
		ctx context.Context,
		symbol string,
		from *model.Date,
		to *model.Date,
		sortDirection model.SortDirection,
	) ([]*model.Candle, error)
}
