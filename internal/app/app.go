package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikitabablo/exchangert/config"
	"github.com/mikitabablo/exchangert/internal/client/crypto"
	"github.com/mikitabablo/exchangert/internal/client/openexchangerates"
	deliveryhttp "github.com/mikitabablo/exchangert/internal/delivery/http"
	"github.com/mikitabablo/exchangert/internal/usecase"
)

type App struct {
	config *config.Config

	engine *gin.Engine
	server *http.Server
}

func NewApp(cfg *config.Config) *App {
	app := &App{
		config: cfg,
	}

	fiatProvider := openexchangerates.NewClient(cfg.OpenExchangeRates.Url, cfg.OpenExchangeRates.AppId)
	fiatUsecase := usecase.NewRatesUsecase(fiatProvider)
	ratesHandler := deliveryhttp.NewRatesHandler(fiatUsecase)

	cryptoProvider := crypto.NewStaticDataProvider()
	cryptoUsecase := usecase.NewCryptoUsecase(cryptoProvider)
	exchangeHandler := deliveryhttp.NewExchangeHandler(cryptoUsecase)

	engine := gin.Default()
	engine.GET("/rates", ratesHandler.GetRates)
	engine.GET("/exchange", exchangeHandler.Exchange)

	app.engine = engine

	return app
}

func (app *App) Run() error {
	app.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", app.config.Server.Port),
		Handler: app.engine,
	}

	if err := app.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (app *App) Stop(ctx context.Context) error {
	return app.server.Shutdown(ctx)
}
