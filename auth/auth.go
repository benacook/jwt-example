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

type fatalCall func(v ...interface{})

var (
	secret []byte
	Fatal fatalCall
)
//======================================================================================

//init gets the config from the config.json file stored on the root.
//it figures out the root folder of the project based on the folder name hard-coded
//here as "jwt-example". this is needed for the tests to work.

func init() {
	Fatal = log.Fatal

	basePath := GetProjectBasePath("jwt-example")

	//read the config file
	configFile := ReadFile(basePath, "/config.json")

	config := Config{}
	_ = json.Unmarshal(configFile, &config)

	//read the file containing the secret
	secret = ReadFile(basePath, config.Secret.Path)
}

//======================================================================================

//reads the specified file from the path and returns its contents

func ReadFile(path, file string) []byte{
	contents, err := ioutil.ReadFile(path + file)
	if err != nil {
		testableFatal(Fatal, err)
		return nil
	}
	return contents
}

//======================================================================================

//gets the base project path from the given project name (i.e. folder name)

func GetProjectBasePath(projectName string) string {

	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	splitStr := strings.SplitAfter(basePath, projectName)
	if len(splitStr) < 2{
		testableFatal(Fatal, errors.New("failed to get bas configFile"))
		return ""
	}

	//loop through the split configFile string to find the last "jwt-example"
	last := 0
	for i:=0 ; i < len(splitStr); i++ {
		if strings.Contains(splitStr[i], projectName){
			last = i
		}
	}

	//construct the base configFile
	basePath = ""
	for i:=0 ; i <= last; i++ {
		basePath += splitStr[i]
	}
	return basePath
}

//======================================================================================

//wrapper for log.fatal to make it testable

func testableFatal(call fatalCall, err error){
	call(err)
}
//======================================================================================

//generates a new token with custom claims and signs it using the using the secret
//from the location defined in the config.json file.

func GenerateToken() (string, error){
	token := NewToken()
	return SignToken(token, secret), nil
}

//======================================================================================

//makes a new token

func NewToken() *jwt.Token{
	return jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"custom": "whatever",
		"nbf": time.Now(),
	})
}

//======================================================================================

//signs a token with the given byte string

func SignToken(token *jwt.Token,secret []byte) string{
	tokenString, err := token.SignedString(secret)
	if err != nil {
		log.Println(err)
		return ""
	}
	return tokenString
}

//======================================================================================

//given a token string, validate it using the secret from the location defined in the
//config.json file. if it is valid print out the nbf claim

func ValidateToken(tokenString string) error{
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{},error){
		//validate signing method is as expected before returning secret
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
		return nil
	}else{
		return errors.New("invalid token")
	}
}

//======================================================================================

//middleware validates a bearer token in a http request,
//if valid calls <next> else returns a 403 status code.

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")

		//no bearer token
		if len(splitToken) < 2{
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte("invalid auth"))
			return
		}
		reqToken = splitToken[1]

		//attempt to validate token
		if err := ValidateToken(reqToken); err == nil {
			//token is valid, run function passed to middleware
			next.ServeHTTP(w, r)
		}else{
			//token invalid, return 403
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte("invalid auth"))
		}
	})
}
