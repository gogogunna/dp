package main

import (
	"context"
	"dp/internal/api/http"
	v1 "dp/internal/api/http/all_methods/v1"
	"dp/internal/client/instrument_info_providers/client"
	"dp/internal/client/operations_provider"
	"dp/internal/client/portfolio_provider"
	"dp/internal/client_factory"
	"dp/internal/main_page"
	"dp/internal/operations/deal_operation_provider"
	portfolio_items_provider "dp/internal/position/portfolio_provider"
	"fmt"
	_ "github.com/hashicorp/go-msgpack/codec"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	cfg := investgo.Config{
		EndPoint:                      "invest-public-api.tinkoff.ru:443",
		Token:                         "",
		AppName:                       "invest-api-go-sdk",
		AccountId:                     "",
		DisableResourceExhaustedRetry: false,
		DisableAllRetry:               false,
		MaxRetries:                    3,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	zapConfig.EncoderConfig.TimeKey = "time"
	l, err := zapConfig.Build()
	logger := l.Sugar()
	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Printf(err.Error())
		}
	}()
	if err != nil {
		log.Fatalf("logger creating error %v", err)
	}

	factory := client_factory.NewClientFactory(cfg, l)
	factory.Start(ctx)

	instrumentsProvider := client.NewInstrumentsInfoProvider()
	portfolioClientProvider := portfolio_provider.NewPortfolioProvider()
	operationsProvider := operations_provider.NewOperationsProvider()
	portfolioProvider := portfolio_items_provider.NewPortfolioProvider(portfolioClientProvider, instrumentsProvider)
	dealOperationsProvider := deal_operation_provider.NewDealOperationsProvider(operationsProvider, instrumentsProvider)
	mainPageInfoProvider := main_page.NewMainPageInfoProvider(portfolioProvider)
	handler := v1.NewHTTPServerHandler(factory, mainPageInfoProvider, portfolioProvider, dealOperationsProvider)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		httpServerErr := http.StartHTTPServer(handler)
		if httpServerErr != nil {
			fmt.Println(httpServerErr)
		}
	}()

	wg.Wait()
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
