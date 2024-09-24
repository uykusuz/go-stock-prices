package ports

import (
	"context"
	"go-stock-prices/model"
)

type NotificationRepository interface {
	NotifyPerformance(ctx context.Context, performance *model.Performance) error
}
