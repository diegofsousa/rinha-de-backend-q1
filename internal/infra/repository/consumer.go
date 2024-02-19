package repository

import (
	"context"
	infraerrors "github.com/diegofsousa/rinha-de-backend-q1/internal/infra"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"time"
)

type Consumer struct {
	database *PgConnection
}

type Statement struct {
	Balance          Balance
	LastTransactions []Transaction
}
type Balance struct {
	Date  time.Time
	Total int64
	Limit int64
}
type Transaction struct {
	Value       int64
	Type        string
	Description string
	CreatedAt   time.Time
}

func NewConsumer(databaseUrl string) Consumer {
	return Consumer{
		database: NewPgConnection(databaseUrl),
	}
}

func (c Consumer) GetLastTransactionsByConsumerId(ctx context.Context, consumerId int64) (*Balance, *[]Transaction, error) {
	conn, err := c.database.Connect(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer c.database.Close(ctx, conn)

	consumerLimits, err := c.getLimitByConsumerId(ctx, conn, consumerId)
	if err != nil {
		return nil, nil, err
	}

	lastTransactions, err := c.getLastTransactionsByConsumerId(ctx, conn, consumerId)
	if err != nil {
		return nil, nil, err
	}

	return consumerLimits, lastTransactions, nil
}

func (c Consumer) InsertNewTransactionByConsumerId(ctx context.Context, consumerId int64, transaction *Transaction) (*Balance, error) {
	conn, err := c.database.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.database.Close(ctx, conn)

	consumerLimits, err := c.getLimitByConsumerId(ctx, conn, consumerId)
	if err != nil {
		return nil, err
	}

	var previewOperation int64

	if transaction.Type == "d" {
		previewOperation = consumerLimits.Total - transaction.Value
		if previewOperation*-1 > consumerLimits.Limit {
			log.Info("unsupported transaction", err)
			return nil, infraerrors.UnsupportedTransaction
		}
	} else {
		previewOperation = consumerLimits.Total + transaction.Value
	}

	err = c.insertTransaction(ctx, conn, consumerId, transaction)
	if err != nil {
		return nil, err
	}

	err = c.updateConsumerBalance(ctx, conn, consumerId, previewOperation)
	if err != nil {
		return nil, err
	}

	return &Balance{
		Total: previewOperation,
		Limit: consumerLimits.Limit,
	}, nil
}

func (c Consumer) getLastTransactionsByConsumerId(ctx context.Context, conn *pgx.Conn, consumerId int64) (*[]Transaction, error) {
	var lastTransactions []Transaction

	rows, err := conn.Query(ctx, "select t.value, t.type, t.description, t.created_at from transactions t where t.consumer_id = $1 order by t.created_at desc limit 10", consumerId)
	if err != nil {
		log.Error("sql error", err)
		return nil, infraerrors.TransactionsErrorSql
	}
	for rows.Next() {
		var l Transaction
		err = rows.Scan(&l.Value, &l.Type, &l.Description, &l.CreatedAt)
		if err != nil {
			log.Error("sql error", err)
			return nil, infraerrors.TransactionsErrorSql
		}
		lastTransactions = append(lastTransactions, l)

	}
	return &lastTransactions, nil
}

func (c Consumer) getLimitByConsumerId(ctx context.Context, conn *pgx.Conn, consumerId int64) (*Balance, error) {
	var b Balance

	row := conn.QueryRow(ctx, "select now(), c.account_limit , c.balance from consumer c where c.id = $1", consumerId)
	err := row.Scan(&b.Date, &b.Limit, &b.Total)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info("not found register")
			return nil, infraerrors.NotFoundConsumerRegister
		}
		log.Error("sql error", err)
		return nil, infraerrors.ConsumerErrorSql
	}

	return &b, nil
}

func (c Consumer) insertTransaction(ctx context.Context, conn *pgx.Conn, consumerId int64, transaction *Transaction) error {
	query := `insert into transactions (consumer_id, value, type, description) values ($1, $2, $3, $4);`
	_, err := conn.Exec(ctx, query, consumerId, transaction.Value, transaction.Type, transaction.Description)
	if err != nil {
		log.Error("insert transaction sql error", err)
		return infraerrors.ConsumerErrorSql
	}

	return nil
}

func (c Consumer) updateConsumerBalance(ctx context.Context, conn *pgx.Conn, consumerId int64, newBalance int64) error {
	query := `update consumer SET balance = $2 WHERE id = $1;`
	command, err := conn.Exec(ctx, query, consumerId, newBalance)
	if err != nil {
		log.Error("update consumer balance sql error", err)
		return infraerrors.ConsumerErrorSql
	}

	if command.RowsAffected() == 0 {
		log.Info("no rows updated")
		return infraerrors.NotFoundConsumerRegister
	}

	return nil
}
