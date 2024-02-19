package gateway

import (
	"context"
	"github.com/diegofsousa/rinha-de-backend-q1/internal/infra/repository"
)

type ConsumerQuery interface {
	GetLastTransactionsByConsumerId(ctx context.Context, consumerId int64) (*repository.Balance, *[]repository.Transaction, error)
	InsertNewTransactionByConsumerId(ctx context.Context, consumerId int64, transaction *repository.Transaction) (*repository.Balance, error)
}
