package middleware

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/vidhill/the-starry-night/domain"
)

type LogMiddleware struct {
	next   http.Handler
	logger domain.LogProvider
}

func (m LogMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	logger := m.logger
	// wrapped writer
	ww := middleware.NewWrapResponseWriter(w, req.ProtoMajor)

	defer func() {
		respStatus := ww.Status()
		s := fmt.Sprintf(`served %s request "%s" returned %v response`, req.Method, req.URL.Path, respStatus)

		if respStatus == http.StatusInternalServerError {
			logger.Error(s)
		} else {
			logger.Info(s)
		}

	}()

	m.next.ServeHTTP(ww, req)
}

func MakeMyLoggerMiddleware(logger domain.LogProvider) func(http.Handler) http.Handler {
	logHandler := handlerFactory(logger)

	return func(next http.Handler) http.Handler {
		return logHandler(next)
	}
}

func handlerFactory(logger domain.LogProvider) func(next http.Handler) LogMiddleware {
	return func(next http.Handler) LogMiddleware {
		return LogMiddleware{
			logger: logger,
			next:   next,
		}
	}
}
