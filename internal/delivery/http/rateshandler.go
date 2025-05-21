package http

import (
    "github.com/mikitabablo/exchangert/internal/domain"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
)

type RatesHandler struct {
    ratesUsecase domain.FiatUsecase
}

func NewRatesHandler(ratesUsecase domain.FiatUsecase) *RatesHandler {
    return &RatesHandler{ratesUsecase: ratesUsecase}
}

func (h *RatesHandler) GetRates(c *gin.Context) {
    ctx := c.Request.Context()

    currencies := c.Query("currencies")
    if currencies == "" {
        c.AbortWithStatus(http.StatusBadRequest)
        return
    }

    list := strings.Split(currencies, ",")
    if len(list) < 2 {
        c.AbortWithStatus(http.StatusBadRequest)
        return
    }

    rates, err := h.ratesUsecase.GetRates(ctx, list)
    if err != nil {
        c.AbortWithStatus(http.StatusBadRequest)
        return
    }

    c.JSON(http.StatusOK, rates)
}
