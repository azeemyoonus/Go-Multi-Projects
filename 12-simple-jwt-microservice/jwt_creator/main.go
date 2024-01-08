package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// MySigninKey is the secret key used to sign the JWT
var MySigninKey = []byte("SECRET_KEY")

func Index(w http.ResponseWriter, r *http.Request) {
	validToken, err := GetJWT()
	fmt.Print(validToken)
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
	fmt.Fprint(w, string(validToken))
}

func GetJWT() ([]byte, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = "amit amrutiya"
	claims["aud"] = "billing.jwtgo.io"
	claims["iss"] = "jwtgo.io"
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()

	tokenString, err := token.SignedString(MySigninKey)
	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return nil, err
	}
	return []byte(tokenString), nil

}
func handleRequests() {
	http.HandleFunc("/", Index)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func main() {
	handleRequests()
}
