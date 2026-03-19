package api

import "net/http"

func (*Config) HandleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("OK"))
	w.WriteHeader(http.StatusOK)
}
