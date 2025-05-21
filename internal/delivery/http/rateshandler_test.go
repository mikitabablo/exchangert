package http_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	delivery "github.com/mikitabablo/exchangert/internal/delivery/http"
	"github.com/mikitabablo/exchangert/internal/domain"
	"github.com/stretchr/testify/assert"
)

type mockFiatUsecase struct{}

func (m *mockFiatUsecase) GetRates(ctx context.Context, currencies []string) ([]domain.FiatRate, error) {
	return []domain.FiatRate{
		{From: "USD", To: "EUR", Rate: 0.9},
		{From: "EUR", To: "USD", Rate: 1.1},
	}, nil
}

type mockBrokenRatesUsecase struct{}

func (m *mockBrokenRatesUsecase) GetRates(ctx context.Context, currencies []string) ([]domain.FiatRate, error) {
	return nil, assert.AnError
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := delivery.NewRatesHandler(&mockFiatUsecase{})
	r.GET("/rates", handler.GetRates)

	return r
}

func TestGetRatesHandler(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/rates?currencies=USD,EUR", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var rates []domain.FiatRate
	err := json.NewDecoder(strings.NewReader(w.Body.String())).Decode(&rates)
	assert.NoError(t, err)
	assert.Len(t, rates, 2)

	assert.Equal(t, "USD", rates[0].From)
	assert.Equal(t, "EUR", rates[0].To)
}

func TestGetRatesHandler_Success(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/rates?currencies=USD,EUR", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var rates []domain.FiatRate
	err := json.NewDecoder(w.Body).Decode(&rates)
	assert.NoError(t, err)
	assert.Len(t, rates, 2)
}

func TestGetRatesHandler_MissingQueryParam(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/rates", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetRatesHandler_OneCurrencyProvided(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/rates?currencies=USD", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetRatesHandler_EmptyCurrencies(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/rates?currencies=", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetRatesHandler_InternalError(t *testing.T) {
	brokenUsecase := &mockBrokenRatesUsecase{}
	handler := delivery.NewRatesHandler(brokenUsecase)

	r := gin.Default()
	r.GET("/rates", handler.GetRates)

	req, _ := http.NewRequest(http.MethodGet, "/rates?currencies=USD,EUR", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
