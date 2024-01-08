package main

import (
	"fmt"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

var MySigninKey = []byte("SECRET_KEY")

func handleRequests() {
	http.HandleFunc("/", isAuthorized(homePage))
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func main() {
	fmt.Println("Hello, Server")
	handleRequests()
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: homePage")
	fmt.Fprintf(w, "Welcome to the HomePage!")
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			fmt.Println("Token: ", r.Header["Token"])
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { // Check the signing method
					return nil, fmt.Errorf("There was an error")
				}
				checkAudience := token.Claims.(jwt.MapClaims)["aud"] == "billing.jwtgo.io"
				checkIssuer := token.Claims.(jwt.MapClaims)["iss"] == "jwtgo.io"
				if !checkAudience || !checkIssuer {
					return nil, fmt.Errorf("There was an error in Claims")
				}
				return MySigninKey, nil
			})
			if err != nil {
				fmt.Println("Error: ", err.Error())
				fmt.Fprint(w, err.Error())
			}
			if token.Valid {
				endpoint(w, r)
			}
		} else {
			fmt.Println("Not Authorized")
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}
