// package main

// import (
// 	"errors"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"
// 	"text/template"
// )

// func initialize() {
// 	var tmpl = template.Must(template.ParseGlob("view/*.html"))
// }

// func getRoot(w http.ResponseWriter, r *http.Request) {

// 	fmt.Printf("got / request\n")
// 	http.FileServer(http.Dir("view/index.html"))
// 	// io.WriteString(w, "Hello, Abdullah!\n")
// }

// func getLogin(w http.ResponseWriter, r *http.Request) {
// 	fmt.Printf("got /login request\n")
// 	io.WriteString(w, "Login, HTTP!\n")
// 	file, err := os.Open("./view/login.html")
// 	if err != nil {
// 		fmt.Fprintf(w, "Error loading file: %s", err)
// 		return
// 	}
// 	defer file.Close()

// 	// Read the file contents
// 	html, err := io.ReadAll(file)
// 	if err != nil {
// 		fmt.Fprintf(w, "Error reading file: %s", err)
// 		return
// 	}

// 	// Write the file contents to the response
// 	w.Header().Set("Content-Type", "text/html")
// 	w.Write(html)
// }

// func getSignup(w http.ResponseWriter, r *http.Request) {
// 	initialize()

// 	a := "Savio"
// 	tmpl.ExecuteTemplate(w, "/signup/", a)

// }

// func main() {

// 	mux := http.NewServeMux()

// 	mux.Handle("/", http.FileServer(http.Dir("./view")))
// 	// mux.HandleFunc("/", getRoot())
// 	mux.HandleFunc("/login", getLogin)

// 	mux.HandleFunc("/signup", getSignup)

// 	// Starting Server in PORT:3333
// 	err := http.ListenAndServe(":3333", mux)
// 	if errors.Is(err, http.ErrServerClosed) {
// 		fmt.Printf("Server Closed")
// 	} else if err != nil {
// 		fmt.Printf("Error starting server: %s", err)
// 	}

// }

package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

var tmpl *template.Template

func getHome(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Home Loaded")
	tmpl.ExecuteTemplate(w, "index.html", "Abdullah Nettoor")
}

func getLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Login Loaded")
	tmpl.ExecuteTemplate(w, "login.html", nil)
}

func getSignup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Signup Loaded")
	tmpl.ExecuteTemplate(w, "signup.html", nil)
}

func postSignup(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "signup.html", nil)
}

func postLogout(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func main() {
	tmpl = template.Must(template.ParseGlob("view/template/*.html"))
	// tmpl, _ = template.ParseGlob("view/template/*.html")

	// Create Server
	fmt.Println("Starting Server...")
	mux := http.NewServeMux()

	http.HandleFunc("/", getHome)

	mux.HandleFunc("/login", getLogin)

	mux.HandleFunc("/signup", getSignup)

	mux.HandleFunc("/logout", postLogout)
	mux.HandleFunc("/signup-post", postSignup)

	// Start Server
	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server Closed")
	} else if err != nil {
		fmt.Printf("Error starting server: %s", err)
	} else {
		fmt.Println("Server Started on PORT:3333")
	}
}
