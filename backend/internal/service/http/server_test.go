package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/durianpay/fullstack-boilerplate/internal/openapigen"
)

type apiStub struct{}

func (apiStub) PostDashboardV1AuthLogin(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (apiStub) GetDashboardV1Payments(w http.ResponseWriter, _ *http.Request, _ openapigen.GetDashboardV1PaymentsParams) {
	w.WriteHeader(http.StatusOK)
}

type tokenVerifierStub struct{ token string }

func (s tokenVerifierStub) VerifyToken(token string) error {
	if token != s.token {
		return http.ErrNoCookie
	}
	return nil
}

func TestServerRequiresAndPassesBearerToken(t *testing.T) {
	server := NewServer(apiStub{}, "../openapi.yaml", tokenVerifierStub{token: "valid-token"})

	unauthorized := httptest.NewRequest(http.MethodGet, "/dashboard/v1/payments", nil)
	unauthorizedResponse := httptest.NewRecorder()
	server.Routes().ServeHTTP(unauthorizedResponse, unauthorized)
	if unauthorizedResponse.Code != http.StatusUnauthorized {
		t.Fatalf("missing token status = %d", unauthorizedResponse.Code)
	}

	authorized := httptest.NewRequest(http.MethodGet, "/dashboard/v1/payments", nil)
	authorized.Header.Set("Authorization", "Bearer valid-token")
	authorizedResponse := httptest.NewRecorder()
	server.Routes().ServeHTTP(authorizedResponse, authorized)
	if authorizedResponse.Code != http.StatusOK {
		t.Fatalf("valid token status = %d, body = %s", authorizedResponse.Code, authorizedResponse.Body.String())
	}
}

func TestServerHealthCheckDoesNotRequireAuthentication(t *testing.T) {
	server := NewServer(apiStub{}, "../openapi.yaml", nil)
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()

	server.Routes().ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("health check status = %d", response.Code)
	}
}

func TestServerServesOpenAPIDocumentation(t *testing.T) {
	server := NewServer(apiStub{}, "../../../../openapi.yaml", nil)
	tests := []struct {
		path         string
		bodyContains string
	}{
		{path: "/openapi.yaml", bodyContains: "openapi: 3.0.3"},
		{path: "/swagger/index.html", bodyContains: "/openapi.yaml"},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.path, nil)
			response := httptest.NewRecorder()

			server.Routes().ServeHTTP(response, request)

			if response.Code != http.StatusOK {
				t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
			}
			if !strings.Contains(response.Body.String(), test.bodyContains) {
				t.Fatalf("body does not contain %q", test.bodyContains)
			}
		})
	}
}
