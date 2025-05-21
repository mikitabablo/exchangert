package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mikitabablo/exchangert/internal/client/crypto"
	delivery "github.com/mikitabablo/exchangert/internal/delivery/http"
	"github.com/mikitabablo/exchangert/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func setupExchangeRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	cryptoUsecase := usecase.NewCryptoUsecase(crypto.NewStaticDataProvider())
	handler := delivery.NewExchangeHandler(cryptoUsecase)

	r.GET("/exchange", handler.Exchange)
	return r
}

func TestExchangeHandler(t *testing.T) {
	router := setupExchangeRouter()

	tests := []struct {
		name           string
		query          string
		expectedStatus int
		expectSuccess  bool
	}{
		{
			name:           "valid request USDT to BEER",
			query:          "/exchange?from=USDT&to=BEER&amount=1.0",
			expectedStatus: http.StatusOK,
			expectSuccess:  true,
		},
		{
			name:           "valid request WBTC to USDT",
			query:          "/exchange?from=WBTC&to=USDT&amount=1.0",
			expectedStatus: http.StatusOK,
			expectSuccess:  true,
		},
		{
			name:           "missing amount param",
			query:          "/exchange?from=USDT&to=BEER",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "negative amount param",
			query:          "/exchange?from=USDT&to=BEER&amount=-1",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "unsupported crypto",
			query:          "/exchange?from=DOGE&to=USDT&amount=1",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tt.query, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectSuccess {
				var resp struct {
					From   string  `json:"from"`
					To     string  `json:"to"`
					Amount float64 `json:"amount"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, 200)
			}
		})
	}
}
