package usecase

import (
	"context"
	"github.com/diegofsousa/rinha-de-backend-q1/internal/application"
	"github.com/diegofsousa/rinha-de-backend-q1/internal/infra/repository"

	"time"
)

type (
	TransactionInput struct {
		Value       int64  `json:"valor"`
		Type        string `json:"tipo"`
		Description string `json:"descricao"`
	}
	TransactionOutput struct {
		Limit   int64 `json:"limite"`
		Balance int64 `json:"saldo"`
	}
)

type (
	StatementsOutput struct {
		Balance          Balance       `json:"saldo"`
		LastTransactions []Transaction `json:"ultimas_transacoes"`
	}
	Balance struct {
		Total int64     `json:"total"`
		Date  time.Time `json:"data_extrato"`
		Limit int64     `json:"limite"`
	}
	Transaction struct {
		Value       int64     `json:"valor"`
		Type        string    `json:"tipo"`
		Description string    `json:"descricao"`
		CreatedAt   time.Time `json:"realizada_em"`
	}
)

type Consumer struct {
	query gateway.ConsumerQuery
}

func NewConsumer(query gateway.ConsumerQuery) Consumer {
	return Consumer{
		query: query,
	}
}

func (c Consumer) Statements(ctx context.Context, consumerId int64) (*StatementsOutput, error) {
	balance, lastTransactions, err := c.query.GetLastTransactionsByConsumerId(ctx, consumerId)
	if err != nil {
		return nil, err
	}

	lt := make([]Transaction, 0, len(*lastTransactions))

	for _, t := range *lastTransactions {
		lt = append(lt, Transaction{
			Value:       t.Value,
			Type:        t.Type,
			Description: t.Description,
			CreatedAt:   t.CreatedAt,
		})
	}
	b := Balance{
		Total: balance.Total,
		Date:  balance.Date,
		Limit: balance.Limit,
	}

	return &StatementsOutput{
		Balance:          b,
		LastTransactions: lt,
	}, nil
}

func (c Consumer) Transaction(ctx context.Context, consumerId int64, input TransactionInput) (*TransactionOutput, error) {
	transaction, err := c.query.InsertNewTransactionByConsumerId(ctx, consumerId, &repository.Transaction{
		Value:       input.Value,
		Type:        input.Type,
		Description: input.Description,
	})
	if err != nil {
		return nil, err
	}

	return &TransactionOutput{
		Balance: transaction.Total,
		Limit:   transaction.Limit,
	}, nil
}
