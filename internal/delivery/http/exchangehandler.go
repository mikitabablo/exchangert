package http

import (
	"github.com/gin-gonic/gin"
	"github.com/mikitabablo/exchangert/internal/domain"
	"net/http"
	"strconv"
)

type ExchangeHandler struct {
	usecase domain.CryptoUsecase
}

func NewExchangeHandler(u domain.CryptoUsecase) *ExchangeHandler {
	return &ExchangeHandler{usecase: u}
}

func (h *ExchangeHandler) Exchange(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	amountStr := c.Query("amount")
	if from == "" || to == "" || amountStr == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	result, err := h.usecase.Exchange(from, to, amount)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, result)
}
