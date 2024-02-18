package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
)

var (
	GenericError = errors.New("generic error")
)

type Clients struct {
}

func NewClients() *Clients {
	return &Clients{}
}

func (api *Clients) Register(server *echo.Echo) {
	server.POST("/clientes/:id/transacoes", api.Transaction)
}

func (api *Clients) Transaction(c echo.Context) error {
	id := c.Param("id")
	fmt.Print(id)

	return c.JSON(http.StatusOK, TransactionResponse{
		Limit:   40,
		Balance: 56,
	})

}

func (api *Clients) handleError(err error, c echo.Context) error {
	switch errors.Cause(err) {
	case GenericError:
		return echo.ErrInternalServerError
	default:
		return echo.ErrInternalServerError
	}

}
