package api

import (
	"net/http"

	"github.com/jadczakd/url-shortener/shortener" 
)

// "github.com/go-chi/chi"
// jsonS "github.com/jadczakd/url-shortener/serializer/json"
// msgPackS "github.com/jadczakd/url-shortener/serializer/msgpack"

type RedirectInterface interface {
	Get(http.ResponseWriter, *http.Request)
	Post(http.ResponseWriter, *http.Request)
}

type handlers struct {
	redirectService shortener.RedirectService
}

func NewHandler(rs shortener.RedirectService) RedirectRepo {
	return &handler{redirectService: rs}
}