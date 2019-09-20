package controllers

import (
	"fmt"
	"net/http"
)

// LoginHandler helps in login and retrieving the user details
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In Login Handler")
}
