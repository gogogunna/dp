package http

import "net/http"

type Handler interface {
	Authenticate(w http.ResponseWriter, r *http.Request)
	MainPageInfo(w http.ResponseWriter, r *http.Request)
}

func StartHTTPServer(handler Handler) error {
	http.HandleFunc("/authenticate/v1", handler.Authenticate)
	http.HandleFunc("/main/v1", handler.MainPageInfo)
	err := http.ListenAndServe(":80", nil)
	return err
}
