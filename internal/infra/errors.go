package infraerrors

import "github.com/pkg/errors"

var (
	GenericError             = errors.New("generic error")
	NotFoundConsumerRegister = errors.New("consumer not found in db")
	ConsumerErrorSql         = errors.New("consumer sql error")
	TransactionsErrorSql     = errors.New("transactions sql error")
	BadRequest               = errors.New("bad request")
	UnsupportedTransaction   = errors.New("unsupported transaction")
)
