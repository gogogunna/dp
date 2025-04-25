package v1

import (
	"context"
	"dp/internal"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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
	Portfolio(
		ctx context.Context,
		accountClient internal.AccountIdWithAttachedClientttt,
		currency int,
	) ([]internal.PortfolioPosition, error)
}

type OperationsProvider interface {
	Operations(
		ctx context.Context,
		accountClient internal.AccountIdWithAttachedClientttt,
		figies []internal.Figi,
		interval internal.TimeInterval,
	) (map[internal.Figi][]internal.EnrichedOperation, error)
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
		DailyPercent:   info.DailyPercent,
		DailyMoney:     info.DailyMoney,
		AlltimePercent: info.AlltimePercent,
		AlltimeMoney:   info.AlltimeMoney,
		AllMoney:       info.AllMoney,
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

	operations, err := h.operationsProvider.Operations(r.Context(), client, figies, interval)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	operationsMapped := make(map[string][]FullOperationInfo, len(operations))
	for figi, figiOperations := range operations {
		infos := make([]FullOperationInfo, 0, len(figiOperations))
		for _, operation := range figiOperations {
			var positionInfo *PositionInfo
			if operation.EnrichingInfo.IsDefined() {
				positionInfo = &PositionInfo{
					Name:     operation.EnrichingInfo.Value().Name,
					LogoPath: operation.EnrichingInfo.Value().LogoPath,
				}
			}
			infos = append(infos, FullOperationInfo{
				Operation: Operation{
					Position: Position{
						Figi:     string(operation.Operation.Position.Figi),
						Price:    operation.Operation.Position.Price,
						Quantity: operation.Operation.Position.Quantity,
					},
					OperationType:        int(operation.Operation.OperationType),
					OperationDescription: operation.Operation.OperationDescription,
				},
				PositionInfo: positionInfo,
			})
		}

		operationsMapped[string(figi)] = infos
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

	info, err := h.portfolioProvider.Portfolio(r.Context(), client, portfolioReq.Currency)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fullPositions := make([]FullPortfolioPositionInfo, 0, len(info))
	for _, position := range info {
		fullPositions = append(fullPositions, FullPortfolioPositionInfo{
			PortfolioPositionInfo: PortfolioPositionInfo{
				Position: Position{
					Figi:     string(position.Position.Figi),
					Price:    position.Position.Price,
					Quantity: position.Position.Quantity,
				},
				AllTimeMoney:   position.Position.AllTimeMoney,
				AllTimePercent: position.Position.AllTimePercent,
				DailyMoney:     position.Position.DailyMoney,
				DailyPercent:   position.Position.DailyPercent,
				AllMoney:       position.Position.AllMoney,
			},
			PositionAdditionalInfo: PositionInfo{
				Name:     position.Enriched.Name,
				LogoPath: position.Enriched.LogoPath,
			},
		})
	}
	resp := PortfolioResponse{
		Positions: fullPositions,
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

type OperationsRequest struct {
	Figies   []string     `json:"figies"`
	Interval TimeInterval `json:"interval"`
}

type Operation struct {
	Position             `json:"position"`
	OperationType        int    `json:"operation_type"`
	OperationDescription string `json:"operation_description"`
}

type FullOperationInfo struct {
	Operation    Operation     `json:"operation"`
	PositionInfo *PositionInfo `json:"position_info"`
}

type OperationsResponse struct {
	Operations map[string][]FullOperationInfo `json:"operations"`
}

type PortfolioRequest struct {
	Currency int `json:"currency"`
}

type PortfolioResponse struct {
	Positions []FullPortfolioPositionInfo `json:"full_positions"`
}

type Position struct {
	Figi     string `json:"figi"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}

type PositionInfo struct {
	Name     string `json:"name"`
	LogoPath string `json:"logo_path"`
}

type PortfolioPositionInfo struct {
	Position       Position `json:"position"`
	AllTimeMoney   int      `json:"all_time_money"`
	AllTimePercent int      `json:"all_time_percent"`
	DailyMoney     int      `json:"daily_money"`
	DailyPercent   int      `json:"daily_percent"`
	AllMoney       int      `json:"all_money"`
}

type FullPortfolioPositionInfo struct {
	PortfolioPositionInfo  PortfolioPositionInfo `json:"portfolio_position_info"`
	PositionAdditionalInfo PositionInfo          `json:"position_additional_info"`
}

type FullPositionInfo struct {
	Position Position     `json:"position"`
	Info     PositionInfo `json:"info"`
}

type TimeInterval struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type MainPageResponse struct {
	Name           string `json:"name"`
	DailyPercent   int    `json:"daily_percent"`
	DailyMoney     int    `json:"daily_money"`
	AlltimePercent int    `json:"alltime_percent"`
	AlltimeMoney   int    `json:"alltime_money"`
	AllMoney       int    `json:"all_money"`
}

type AuthResponse struct {
	OK string `json:"ok"`
}
