package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/taaaaakahiro/golang-rest-example/pkg/config"
)

func CORSHeaderMiddleware(cfg *config.Config) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			allowOrigin := cfg.Server.AllowCorsOrigin

			w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			//w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")  // <- CORSMethodMiddleware でつけてる。
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			next.ServeHTTP(w, req)
		})
	}
}
