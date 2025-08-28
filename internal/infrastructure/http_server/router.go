package http_server

import (
	"net/http"
	"webhook-processor-ms/internal/infrastructure/interceptor"

	"github.com/gorilla/mux"
)

// NewRouter configura e retorna um novo roteador HTTP.
// Ele recebe os handlers por injeção de dependência.
func NewRouters(routeHandles []RouterInterface) *mux.Router {
	router := mux.NewRouter()

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			interceptor.CorsMiddleware(w, r)
			next.ServeHTTP(w, r)
		})
	})

	for _, route := range routeHandles {
		registerRoute(route, router)
	}

	return router
}

func registerRoute(routeHandle RouterInterface, router *mux.Router) {
	for _, model := range routeHandle.getRoutersModel() {

		if model.HasAuthenticated {
			router.HandleFunc(model.URI, func(w http.ResponseWriter, r *http.Request) {}).Methods(http.MethodOptions)
			router.HandleFunc(model.URI, interceptor.Autenticar(model.Func)).Methods(model.Method)
		} else {
			router.HandleFunc(model.URI, func(w http.ResponseWriter, r *http.Request) {}).Methods(http.MethodOptions)
			router.HandleFunc(model.URI, interceptor.Logger(model.Func)).Methods(model.Method)
		}
	}
}
