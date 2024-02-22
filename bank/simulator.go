package bank

import "context"

type Simulator interface {
	StartPaymentProcess(ctx context.Context) (string, error)
	RefundPayment(ctx context.Context)
}
