package main

import (
	"context"
	"dp/internal"
	"dp/internal/api/http"
	v1 "dp/internal/api/http/all_methods/v1"
	"dp/internal/client_factory"
	"dp/internal/main_page"
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
	fmt.Println(internal.OperationTypeInputSwift)
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

	mainPageInfoProvider := main_page.NewMainPageInfoProvider()

	handler := v1.NewHTTPServerHandler(factory, mainPageInfoProvider)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		httpServerErr := http.StartHTTPServer(handler)
		if httpServerErr != nil {
			logger.Fatalf("http server died" + err.Error())
		}
	}()

	wg.Wait()
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
