package middleware

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/vidhill/the-starry-night/service"
)

type LogMiddleware struct {
	next   http.Handler
	logger service.LoggerService
}

func (m LogMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	// wrapped writer
	ww := middleware.NewWrapResponseWriter(w, req.ProtoMajor)

	defer func() {
		s := fmt.Sprintf(`served %s request "%s" returned %v response`, req.Method, req.URL.Path, ww.Status())
		m.logger.Info(s)
	}()

	m.next.ServeHTTP(ww, req)
}

func MakeMyLoggerMiddleware(logger service.LoggerService) func(http.Handler) http.Handler {
	logHandler := handlerFactory(logger)

	return func(next http.Handler) http.Handler {
		return logHandler(next)
	}
}

func handlerFactory(logger service.LoggerService) func(next http.Handler) LogMiddleware {
	return func(next http.Handler) LogMiddleware {
		return LogMiddleware{
			logger: logger,
			next:   next,
		}
	}
}
