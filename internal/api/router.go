package api

import (
	"net/http"

	_ "github.com/diogocarasco/go-pharmacy-service/docs"
	"github.com/diogocarasco/go-pharmacy-service/internal/auth"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

type RouterConfig struct {
	Handlers      *Handlers
	Authenticator *auth.Authenticator
}

func NewRouter(cfg RouterConfig) *mux.Router {
	r := mux.NewRouter()

	r.Use(MetricsMiddleware)

	r.HandleFunc("/health", cfg.Handlers.HealthCheckHandler).Methods("GET")
	r.HandleFunc("/metrics", promhttp.Handler().ServeHTTP).Methods("GET")

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	authRouter := r.PathPrefix("/").Subrouter()
	authRouter.Use(cfg.Authenticator.AuthMiddleware)

	authRouter.HandleFunc("/claim", cfg.Handlers.SubmitClaimHandler).Methods("POST")
	authRouter.HandleFunc("/claim/{id}", cfg.Handlers.GetClaimByIDHandler).Methods("GET")
	authRouter.HandleFunc("/reversal", cfg.Handlers.ReverseClaimHandler).Methods("POST")

	return r
}
