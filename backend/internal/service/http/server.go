package http

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/durianpay/fullstack-boilerplate/internal/openapigen"
	"github.com/durianpay/fullstack-boilerplate/internal/transport"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	oapinethttpmw "github.com/oapi-codegen/nethttp-middleware"
)

type Server struct {
	router http.Handler
}

type TokenVerifier interface {
	VerifyToken(rawToken string) error
}

const (
	readTimeout  = 10
	writeTimeout = 10
	idleTimeout  = 60
)

func NewServer(apiHandler openapigen.ServerInterface, openapiYamlPath string, tokenVerifier TokenVerifier) *Server {
	swagger, err := openapigen.GetSwagger()
	if err != nil {
		log.Fatalf("failed to load swagger: %v", err)
	}

	r := chi.NewRouter()

	r.Route("/", func(api chi.Router) {
		api.Use(oapinethttpmw.OapiRequestValidatorWithOptions(
			swagger,
			&oapinethttpmw.Options{
				Options: openapi3filter.Options{
					AuthenticationFunc: bearerAuthentication(tokenVerifier),
				},
				DoNotValidateServers:  true,
				SilenceServersWarning: true,
				ErrorHandler: func(w http.ResponseWriter, _ string, statusCode int) {
					switch statusCode {
					case http.StatusUnauthorized:
						transport.WriteJSONError(w, statusCode, "missing or invalid token")
					case http.StatusBadRequest:
						transport.WriteJSONError(w, statusCode, "invalid request")
					default:
						transport.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
					}
				},
			},
		))
		openapigen.HandlerFromMux(apiHandler, api)
	})

	return &Server{
		router: r,
	}
}

func bearerAuthentication(tokenVerifier TokenVerifier) openapi3filter.AuthenticationFunc {
	return func(_ context.Context, input *openapi3filter.AuthenticationInput) error {
		if input.SecuritySchemeName != "bearerAuth" || tokenVerifier == nil {
			return http.ErrNoCookie
		}

		authorization := input.RequestValidationInput.Request.Header.Get("Authorization")
		const prefix = "Bearer "
		if !strings.HasPrefix(authorization, prefix) {
			return http.ErrNoCookie
		}

		return tokenVerifier.VerifyToken(strings.TrimPrefix(authorization, prefix))
	}
}

func (s *Server) Start(addr string) {
	service := &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  readTimeout * time.Second,
		WriteTimeout: writeTimeout * time.Second,
		IdleTimeout:  idleTimeout * time.Second,
	}
	go func() {
		log.Printf("listening on %s", addr)
		if err := service.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server stopped unexpectedly: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down gracefully...")

	// Timeout for shutdown
	const shutdownTimeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := service.Shutdown(ctx); err != nil {
		log.Fatalf("Forced shutdown: %v", err)
	}

	log.Println("Server stopped cleanly ✔")
}

func (s *Server) Routes() http.Handler {
	return s.router
}
