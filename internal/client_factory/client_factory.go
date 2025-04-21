package client_factory

import (
	"context"
	"dp/internal"
	internal_errors "dp/internal/errors"
	"errors"
	"fmt"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"sync"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	"go.uber.org/zap"
	"sync/atomic"
)

type ClientFactory struct {
	config    investgo.Config
	clients   map[internal.UserToken]internal.AccountIdWithAttachedClientttt
	isStarted atomic.Bool
	logger    *zap.Logger
	mu        sync.RWMutex
	appCtx    context.Context
}

func NewClientFactory(config investgo.Config, logger *zap.Logger) *ClientFactory {
	return &ClientFactory{
		config:  config,
		logger:  logger,
		clients: make(map[internal.UserToken]internal.AccountIdWithAttachedClientttt, 10),
	}
}

func (f *ClientFactory) Start(ctx context.Context) {
	if f.isStarted.CompareAndSwap(false, true) {
		f.appCtx = ctx
		go func() {
			<-ctx.Done()
			for _, userClient := range f.clients {
				err := userClient.Client.Stop()
				if err != nil {
					f.logger.Error("failed to stop userClient", zap.Error(err))
				}
			}
		}()
	}
}

func (f *ClientFactory) Client(ctx context.Context, token internal.UserToken) (internal.AccountIdWithAttachedClientttt, error) {
	var accountClient internal.AccountIdWithAttachedClientttt
	if !f.isStarted.Load() {
		return accountClient, errors.New("client factory is not started")
	}

	if func() bool {
		f.mu.RLock()
		defer f.mu.RUnlock()
		var ok bool
		accountClient, ok = f.clients[token]
		return ok
	}() {
		return accountClient, nil
	}

	cfg := f.config
	cfg.Token = token

	client, err := investgo.NewClient(f.appCtx, cfg, f.logger.Sugar())
	if err != nil {
		return accountClient, fmt.Errorf("failed to create client: %w", err)
	}

	status := pb.AccountStatus_ACCOUNT_STATUS_OPEN
	resp, err := client.NewUsersServiceClient().GetAccounts(&status)
	if err != nil {
		return accountClient, fmt.Errorf("failed to get accounts: %w", err)
	}

	if len(resp.GetAccounts()) == 0 {
		return accountClient, internal_errors.NotFoundAccount
	}
	accountId := resp.GetAccounts()[0].GetId()

	accountClient = internal.AccountIdWithAttachedClientttt{
		AccountId: accountId,
		Client:    client,
	}

	f.mu.Lock()
	f.clients[token] = accountClient
	f.mu.Unlock()

	return accountClient, nil
}
