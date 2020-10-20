package handlers

import (
	"github.com/benacook/jwt-example/auth"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewPublicHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/new-token",
		strings.NewReader(""))
	if err != nil {
		log.Fatal(err)
	}

	rec := httptest.NewRecorder()

	gth := NewGenerateTokenHandler()
	handler := http.HandlerFunc(gth.ServeHTTP)

	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK{
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	splitStr := strings.Split(rec.Body.String(), "token=")

	if len(splitStr) < 2{
		t.Error("did not receive token")
		return
	}

	token := splitStr[1]
	if token == ""{
		t.Error("token is empty")
	}

}

func TestNewRestrictedHandlerFail(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/secret-area",
		strings.NewReader(""))
	if err != nil {
		log.Fatal(err)
	}

	rec := httptest.NewRecorder()

	rh := NewRestrictedHandler()
	handler := http.HandlerFunc(auth.Middleware(rh).ServeHTTP)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusForbidden{
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}

	if rec.Body.String() != "invalid auth"{
		t.Errorf("expected failed auth, got: %v", rec.Body.String())
	}
}

func TestNewRestrictedHandlerPass(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/secret-area",
		strings.NewReader(""))
	if err != nil {
		log.Fatal(err)
	}
	token, _ := auth.GenerateToken()
	req.Header.Add("Authorization", "Bearer " + token)

	rec := httptest.NewRecorder()

	rh := NewRestrictedHandler()
	handler := auth.Middleware(rh)

	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK{
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rec.Body.String() != "hello from the secret area!"{
		t.Errorf("failed to Auth, got: %v", rec.Body.String())
	}
}


