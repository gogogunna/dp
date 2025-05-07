package http

import "net/http"

type Handler interface {
	Authenticate(w http.ResponseWriter, r *http.Request)
	MainPageInfo(w http.ResponseWriter, r *http.Request)
	Portfolio(w http.ResponseWriter, r *http.Request)
	Operations(w http.ResponseWriter, r *http.Request)
}

func StartHTTPServer(handler Handler) error {
	http.HandleFunc("/authenticate/v1", handler.Authenticate) // переделать на запихивание клиента в контекст
	http.HandleFunc("/main/v1", handler.MainPageInfo)
	http.HandleFunc("/portfolio/v1", handler.Portfolio)
	http.HandleFunc("/operations/v1", handler.Operations)
	err := http.ListenAndServeTLS(
		":443",
		"/etc/ssl/certs/fullchain1.pem",
		"/etc/ssl/certs/privkey1.pem",
		nil,
	)
	return err
}
