package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name     string
	Email    string
	Password string
}

var Users map[string]User

var user User

var tmpl *template.Template

func hashPassword(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed)
}

func checkPassword(password string) bool {
	v := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if v == nil {
		return true
	} else {
		return false
	}
}

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
	r.ParseForm()

	user.Name = r.FormValue("name")
	user.Email = r.Form.Get("email")
	// user, ok := Users[r.Form.Get("email")]
	// if ok {
	// 	tmpl.ExecuteTemplate(w, "signup.html", "User already exist, Use another email")
	// 	return
	// }
	hashedPwd := hashPassword(r.Form.Get("password"))
	user.Password = hashedPwd
	if user.Name == "" || user.Email == "" {
		tmpl.ExecuteTemplate(w, "signup.html", nil)
	} else {
		Users[user.Email] = user
		fmt.Println(Users)
	}

	fmt.Println(user)

	tmpl.ExecuteTemplate(w, "index.html", user)
}

func postLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	user, ok := Users[r.FormValue("email")]
	if ok {
		if user.Email != r.FormValue("email") {
			tmpl.ExecuteTemplate(w, "login.html", "Enter a valid email")
		} else if !checkPassword(r.FormValue("password")) {
			tmpl.ExecuteTemplate(w, "login.html", "Incorrect Password")
			return
		}
		tmpl.ExecuteTemplate(w, "index.html", user)
	} else {
		tmpl.ExecuteTemplate(w, "login.html", "User don't exist")
	}
}

func postLogout(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func main() {
	tmpl = template.Must(template.ParseGlob("view/template/*.html"))
	Users = make(map[string]User)

	// Create Server
	fmt.Println("Starting Server...")
	mux := http.NewServeMux()

	// Get Static files like CSS, Images etc...
	mux.Handle("/view/static/", http.StripPrefix("/view/static/", http.FileServer(http.Dir("view/static/"))))

	mux.HandleFunc("/", getHome)
	mux.HandleFunc("/login", getLogin)
	mux.HandleFunc("/login-post", postLogin)
	mux.HandleFunc("/signup", getSignup)
	mux.HandleFunc("/signup-post", postSignup)
	mux.HandleFunc("/logout", postLogout)

	// Start Server
	fmt.Println("Server Started on PORT:3333")
	err := http.ListenAndServe(":3333", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server Closed")
	} else if err != nil {
		fmt.Printf("Error starting server: %s", err)
	}
}
