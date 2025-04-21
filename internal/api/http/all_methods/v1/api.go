package v1

import (
	"context"
	"dp/internal"
	"encoding/json"
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

type HTTPServerHandler struct {
	clientFactory        ClientFactory
	mainPageInfoProvider MainPageInfoProvider
}

func NewHTTPServerHandler(clientFactory ClientFactory, provider MainPageInfoProvider) *HTTPServerHandler {
	return &HTTPServerHandler{
		clientFactory:        clientFactory,
		mainPageInfoProvider: provider,
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
		Name:    info.UserName,
		Percent: int(info.DailyPercent),
		Money:   int(info.DailyMoney),
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

type MainPageResponse struct {
	Name    string `json:"name"`
	Percent int    `json:"daily_percent"`
	Money   int    `json:"daily_money"`
}

type AuthResponse struct {
	OK string `json:"ok"`
}
