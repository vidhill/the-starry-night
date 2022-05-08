package handlers

import "net/http"

var okMessage = []byte("ok")

type Handlers struct {
	// service : FooService
}

type ErrorResponse struct {
	Message   string `json:"message"`
	ErrorCode int    `json:"error_code"`
}

func (s Handlers) Health(w http.ResponseWriter, req *http.Request) {
	w.Write(okMessage)
}

func (s Handlers) GetFoo(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("foo"))
}

func (h Handlers) Hello(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{ "World": true }`))
}

func NewHandlers() Handlers {
	return Handlers{}
}
