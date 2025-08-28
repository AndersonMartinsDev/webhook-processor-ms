package interceptor

import (
	"fmt"
	"log/slog"
	"net/http"
	"webhook-processor-ms/internal/infrastructure/configuration"
)

// CorsMiddleware define os headers CORS
func CorsMiddleware(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", configuration.Origin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, User-Agent")
	w.Header().Set("Access-Control-Allow-Credentials", "false")

	slog.Info(fmt.Sprintf("%s %s %s", r.Method, r.RequestURI, r.Host))

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

}

// Logger escreve informações da requisição no terminal
func Logger(proximaFuncao http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proximaFuncao(w, r)
	}
}

// Autenticar verifica se o usuário fazendo a requisição está autenticado
func Autenticar(proximaFuncao http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// handleToken := security.NewHandlerToken()
		// if erro := handleToken.ValidarToken(r); erro != nil {
		// 	response.Erro(w, http.StatusUnauthorized, erro)
		// 	return
		// }
		proximaFuncao(w, r)
	}
}
