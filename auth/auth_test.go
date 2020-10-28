package auth

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func init() {
	Fatal = func(v ...interface{}){}
}

func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken()
	if err != nil {
		t.Fatal(err)
	}
	if token == ""{
		t.Fatal("received empty token")
	}
}

func TestValidateToken(t *testing.T) {
	token, _ := GenerateToken()
	err := ValidateToken(token)
	if err != nil {
		t.Fatal(err)
	}
}

func TestValidateTokenFail(t *testing.T) {
	err := ValidateToken("an invalid token")
	if err != nil {
		t.Log(err)
	}else{
		t.Fatal("expected error from using an invalid token")
	}
}

func TestReadFileFail(t *testing.T) {
	if data := ReadFile("some/unknown/path", "non-existent-file.nope"); data != nil{
		t.Fatalf("expected empty string, got: %v", data)
	}


}

func TestGetProjectBasePathFail(t *testing.T) {
	path := GetProjectBasePath("IDontExist")
	if path != ""{
		t.Fatalf("expected empty string, got: %v", path)
	}
}

func TestMiddlewareValidAuth(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/secret-area",
		strings.NewReader(""))
	if err != nil {
		log.Fatal(err)
	}
	token, _ := GenerateToken()
	req.Header.Add("Authorization", "Bearer " + token)

	rec := httptest.NewRecorder()

	hh := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := http.HandlerFunc(Middleware(hh).ServeHTTP)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK{
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rec.Body.String() == "invalid auth"{
		t.Errorf("failed auth, got: %v", rec.Body.String())
	}
}

func TestMiddlewareInvalidAuth(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/secret-area",
		strings.NewReader(""))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer InvalidToken")

	rec := httptest.NewRecorder()

	hh := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := http.HandlerFunc(Middleware(hh).ServeHTTP)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusForbidden{
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}

	if rec.Body.String() != "invalid auth"{
		t.Errorf("expected failed auth, got: %v", rec.Body.String())
	}
}


func TestMiddlewareNoToken(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/secret-area",
		strings.NewReader(""))
	if err != nil {
		log.Fatal(err)
	}
	rec := httptest.NewRecorder()

	hh := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := http.HandlerFunc(Middleware(hh).ServeHTTP)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusForbidden{
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}

	if rec.Body.String() != "invalid auth"{
		t.Errorf("expected failed auth, got: %v", rec.Body.String())
	}
}
