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

	// Create Server
	fmt.Println("Starting Server...")
	mux := http.NewServeMux()

	// Get Static files like CSS, Images etc...
	mux.Handle("/view/static/", http.StripPrefix("/view/static/", http.FileServer(http.Dir("view/static/"))))

	mux.HandleFunc("/", getHome)
	mux.HandleFunc("/login", getLogin)
	mux.HandleFunc("/signup", getSignup)
	mux.HandleFunc("/logout", postLogout)
	mux.HandleFunc("/signup-post", postSignup)

	// Start Server
	fmt.Println("Server Started on PORT:3333")
	err := http.ListenAndServe(":3333", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server Closed")
	} else if err != nil {
		fmt.Printf("Error starting server: %s", err)
	}
}
