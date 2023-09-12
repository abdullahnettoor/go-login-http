package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	http.FileServer(http.Dir("view/index.html"))
	// io.WriteString(w, "Hello, Abdullah!\n")
}

func getLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /login request\n")
	io.WriteString(w, "Login, HTTP!\n")
	file, err := os.Open("./view/login.html")
	if err != nil {
		fmt.Fprintf(w, "Error loading file: %s", err)
		return
	}
	defer file.Close()

	// Read the file contents
	html, err := io.ReadAll(file)
	if err != nil {
		fmt.Fprintf(w, "Error reading file: %s", err)
		return
	}

	// Write the file contents to the response
	w.Header().Set("Content-Type", "text/html")
	w.Write(html)
}

func main() {
	// tmpl := template.Must(template.ParseGlob("*.html"))

	// tmpl.ExecuteTemplate()

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./view")))
	// mux.HandleFunc("/", getRoot())
	mux.HandleFunc("/login", getLogin)

	err := http.ListenAndServe(":3333", mux)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server Closed")
	} else if err != nil {
		fmt.Printf("Error starting server: %s", err)
	}

}
