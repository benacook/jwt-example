package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Config struct {
	Secret struct{
		Path string `json:"path"`
	}`json:"secret"`
}

var (
	secret []byte
)

func init() {
	var err error
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	splitStr := strings.SplitAfter(basePath, "src")
	if len(splitStr) < 1{
		log.Fatal("failed to get bas path")
	}

	last := 0
	for i:=0 ; i < len(splitStr); i++ {
		if strings.Contains(splitStr[i], "src"){
			last = i
		}
	}

	basePath = ""
	for i:=0 ; i <= last; i++ {
		basePath += splitStr[i]
	}
	configPath := basePath + "/config.json"
	path, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	config := Config{}
	_ = json.Unmarshal(path, &config)

	secret, err = ioutil.ReadFile(basePath + config.Secret.Path)
	if err != nil {
		log.Fatal(err)
	}
}

func GenerateToken() (string, error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"custom": "whatever",
		"nbf": time.Now(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(tokenString string) error{
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{},error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		log.Println(err)
		return err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("valid")
		fmt.Println(claims["nbf"])
	}
	return nil
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		if len(splitToken) < 2{
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte("invalid auth"))
			return
		}
		reqToken = splitToken[1]

		if err := ValidateToken(reqToken); err == nil {
			next.ServeHTTP(w, r)
		}else{
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte("invalid auth"))
		}
	})
}
