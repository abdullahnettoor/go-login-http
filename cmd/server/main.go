package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
}
func getLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /login request\n")
	io.WriteString(w, "Login, HTTP!\n")
}

func main() {

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/login", getLogin)

	err := http.ListenAndServe(":3333", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server Closed")
	} else if err != nil {
		fmt.Printf("Error starting server: %s", err)
	}

}
