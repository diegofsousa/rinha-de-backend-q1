package api

import (
	"github.com/diegofsousa/rinha-de-backend-q1/internal/application/usecase"
	infraerrors "github.com/diegofsousa/rinha-de-backend-q1/internal/infra"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type Consumer struct {
	usecase usecase.Consumer
}

func NewConsumer(usecase usecase.Consumer) *Consumer {
	return &Consumer{
		usecase: usecase,
	}
}

func (api *Consumer) Register(server *echo.Echo) {
	server.POST("/clientes/:id/transacoes", api.PostTransaction)
	server.GET("/clientes/:id/extrato", api.GetTransactions)
}

func (api *Consumer) PostTransaction(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		log.Error("consumer id are required")
		return api.handleError(infraerrors.BadRequest)
	}
	idParsed, err := strconv.Atoi(id)
	if err != nil {
		log.Error("consumer id invalid")
		return api.handleError(infraerrors.BadRequest)
	}
	var request usecase.TransactionInput
	if err := c.Bind(&request); err != nil {
		log.Error("errors in required fields")
		return api.handleError(infraerrors.BadRequest)
	}

	if request.Value <= 0 || request.Description == "" || !(request.Type == "c" || request.Type == "d") {
		log.Error("errors in required fields")
		return api.handleError(infraerrors.BadRequest)
	}

	ctx := c.Request().Context()

	transactions, err := api.usecase.Transaction(ctx, int64(idParsed), request)

	if err != nil {
		return api.handleError(err)
	}

	return c.JSON(http.StatusOK, transactions)
}

func (api *Consumer) GetTransactions(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		log.Error("consumer id are required")
		return api.handleError(infraerrors.BadRequest)
	}
	idParsed, err := strconv.Atoi(id)
	if err != nil {
		log.Error("consumer id invalid")
		return api.handleError(infraerrors.BadRequest)
	}

	ctx := c.Request().Context()

	transactions, err := api.usecase.Statements(ctx, int64(idParsed))

	if err != nil {
		return api.handleError(err)
	}

	return c.JSON(http.StatusOK, transactions)
}

func (api *Consumer) handleError(err error) error {
	switch errors.Cause(err) {
	case infraerrors.NotFoundConsumerRegister:
		return echo.ErrNotFound
	case infraerrors.UnsupportedTransaction:
		return echo.ErrUnprocessableEntity
	case infraerrors.BadRequest:
		return echo.ErrBadRequest
	default:
		return echo.ErrInternalServerError
	}

}
