package v1

import (
	"context"
	"dp/internal"
	"dp/pkg/slices"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type MainPageInfoProvider interface {
	MainPageInfo(
		ctx context.Context,
		accountClient internal.AccountIdWithAttachedClientttt,
	) (internal.MainPageInfo, error)
}

type ClientFactory interface {
	Client(ctx context.Context, token internal.UserToken) (internal.AccountIdWithAttachedClientttt, error)
}

type PortfolioProvider interface {
	PortfolioItems(
		ctx context.Context,
		accountClient internal.AccountIdWithAttachedClientttt,
		currency int,
	) (internal.Portfolio, error)
}

type OperationsProvider interface {
	DealOperations(
		ctx context.Context,
		accountClient internal.AccountIdWithAttachedClientttt,
		figies []internal.Figi,
		interval internal.TimeInterval,
	) (map[internal.Figi][]internal.DealOperation, error)
}

type HTTPServerHandler struct {
	clientFactory        ClientFactory
	mainPageInfoProvider MainPageInfoProvider
	portfolioProvider    PortfolioProvider
	operationsProvider   OperationsProvider
}

func NewHTTPServerHandler(
	clientFactory ClientFactory,
	provider MainPageInfoProvider,
	portfolioProvider PortfolioProvider,
	operationsProvider OperationsProvider,
) *HTTPServerHandler {
	return &HTTPServerHandler{
		clientFactory:        clientFactory,
		mainPageInfoProvider: provider,
		portfolioProvider:    portfolioProvider,
		operationsProvider:   operationsProvider,
	}
}

func (h *HTTPServerHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	_, err := h.clientFactory.Client(r.Context(), token)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	resp := AuthResponse{OK: "OK"}
	bytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (h *HTTPServerHandler) MainPageInfo(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	client, err := h.clientFactory.Client(r.Context(), token)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	info, err := h.mainPageInfoProvider.MainPageInfo(r.Context(), client)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	resp := MainPageResponse{
		Name:           info.UserName,
		DailyPercent:   int(info.PortfolioAnalytics.DailyPercent),
		DailyMoney:     int(info.PortfolioAnalytics.DailyMoney),
		AlltimePercent: int(info.PortfolioAnalytics.AlltimePercent),
		AlltimeMoney:   int(info.PortfolioAnalytics.AlltimeMoney),
		AllMoney:       int(info.PortfolioAnalytics.AllMoney),
	}
	bytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (h *HTTPServerHandler) Operations(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	client, err := h.clientFactory.Client(r.Context(), token)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to read body: %w", err).Error(), 500)
	}

	operationsReq := OperationsRequest{}
	err = json.Unmarshal(bytes, &operationsReq)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to unmarshal body: %w", err).Error(), 500)
		return
	}

	figies := make([]internal.Figi, 0, len(operationsReq.Figies))
	for _, figi := range operationsReq.Figies {
		figies = append(figies, internal.Figi(figi))
	}

	interval := internal.TimeInterval{
		Start: operationsReq.Interval.From,
		End:   operationsReq.Interval.To,
	}

	operations, err := h.operationsProvider.DealOperations(r.Context(), client, figies, interval)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	operationsMapped := make(map[string][]DealOperation, len(operations))
	for figi, dealOperations := range operations {
		operationsMapped[string(figi)] = slices.Convert(dealOperations, mapDealOperation)
	}
	resp := OperationsResponse{
		Operations: operationsMapped,
	}
	bytes, err = json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (h *HTTPServerHandler) Portfolio(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	client, err := h.clientFactory.Client(r.Context(), token)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to read body: %w", err).Error(), 500)
	}

	portfolioReq := PortfolioRequest{}
	err = json.Unmarshal(bytes, &portfolioReq)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to unmarshal body: %w", err).Error(), 500)
		return
	}

	portfolio, err := h.portfolioProvider.PortfolioItems(r.Context(), client, portfolioReq.Currency)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fullPositions := slices.Convert(portfolio.Items, mapPortfolioItem)
	resp := PortfolioResponse{
		Items: fullPositions,
	}
	if portfolio.WarningMessage.IsDefined() {
		message := portfolio.WarningMessage.Value()
		resp.WarningMessage = &message
	}

	bytes, err = json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
