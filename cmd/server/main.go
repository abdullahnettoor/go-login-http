package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User Model
type User struct {
	Name     string
	Email    string
	Password string
}

// Created Users map to store multiple users
var Users map[string]User

var user User

var tmpl *template.Template

// Session Model
type session struct {
	username string
	expiry   time.Time
}

func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

// Created Sessions map to store multiple sessions
var sessions = map[string]session{}

// Function to Create Session
func createSession(u User) (sessionToken string, expiresAt time.Time) {
	sessionToken = uuid.NewString()
	expiresAt = time.Now().Add(time.Hour)
	fmt.Println("\nSession :", sessions)
	sessions[sessionToken] = session{
		username: u.Name,
		expiry:   expiresAt,
	}
	fmt.Println("\nSession Created:", sessions[sessionToken])
	return
}

// Function to Clear Browser Cache
func clearCache(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, no-transform, must-revalidate, private, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("X-Accel-Expires", "0")
}

// Function to generate salted password
func hashPassword(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed)

}

// Function to compare the password with hashed on
func comparePassword(password string, user User) bool {
	v := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if v == nil {
		return true
	} else {
		return false
	}
}

// Handler Function to load home
func getHome(w http.ResponseWriter, r *http.Request) {
	clearCache(w, r)

	// Check if cookie exists
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
		}
		tmpl.ExecuteTemplate(w, "login.html", "Login to see home")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionToken := c.Value

	userSession, exists := sessions[sessionToken]
	if !exists {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

	if userSession.isExpired() {
		delete(sessions, sessionToken)
		tmpl.ExecuteTemplate(w, "login.html", "Login expired")
		return
	} else {
		fmt.Println("\nGET Home Loaded", user)
		tmpl.ExecuteTemplate(w, "index.html", user)
		return
	}

}

// Handler Function to load login
func getLogin(w http.ResponseWriter, r *http.Request) {
	clearCache(w, r)

	// Check if cookie exists
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println(err)
		}
		tmpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

	sessionToken := c.Value

	userSession, exists := sessions[sessionToken]
	if !exists {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		w.WriteHeader(http.StatusUnauthorized)
	}

	if userSession.isExpired() {
		delete(sessions, sessionToken)
		tmpl.ExecuteTemplate(w, "login.html", "Login expired")
		return
	} else {
		fmt.Println("\nHome : GET Login")
		tmpl.ExecuteTemplate(w, "index.html", user)
	}
}

// Handler Function to load Signup
func getSignup(w http.ResponseWriter, r *http.Request) {
	clearCache(w, r)

	// Check if cookie exists
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
		}
		tmpl.ExecuteTemplate(w, "signup.html", nil)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionToken := c.Value

	userSession, exists := sessions[sessionToken]
	if !exists {
		fmt.Println("There is no Session, Signup Loaded")
		tmpl.ExecuteTemplate(w, "signup.html", nil)
		w.WriteHeader(http.StatusUnauthorized)
	}

	if userSession.isExpired() {
		delete(sessions, sessionToken)
		tmpl.ExecuteTemplate(w, "signup.html", nil)
		return
	} else {
		fmt.Println("\nHome : GET Signup", user)
		tmpl.ExecuteTemplate(w, "index.html", user)
	}

}

// Handler Function to Create user
func postSignup(w http.ResponseWriter, r *http.Request) {
	clearCache(w, r)

	r.ParseForm()

	user.Name = r.FormValue("name")
	user.Email = r.Form.Get("email")
	_, ok := Users[r.Form.Get("email")]
	if ok {
		tmpl.ExecuteTemplate(w, "signup.html", "User already exist, Use another email")
		return
	}
	hashedPwd := hashPassword(r.Form.Get("password"))
	user.Password = hashedPwd
	if user.Name == "" || user.Email == "" {
		tmpl.ExecuteTemplate(w, "signup.html", nil)
	} else {
		Users[user.Email] = user
		fmt.Println(Users)

		// Session and Cookie created
		token, expiry := createSession(Users[user.Email])
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   token,
			Expires: expiry,
		})

		fmt.Println("\nHome : POST Signup", user)
		tmpl.ExecuteTemplate(w, "index.html", user)
	}

}

// Handler Function to validate user when login
func postLogin(w http.ResponseWriter, r *http.Request) {
	clearCache(w, r)

	_, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
		}
		tmpl.ExecuteTemplate(w, "login.html", nil)
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	// Retreive data from form
	r.ParseForm()

	// Check if user with entered email exist
	user, exist := Users[r.FormValue("email")]
	if exist {
		if user.Email != r.FormValue("email") {
			tmpl.ExecuteTemplate(w, "login.html", "Enter a valid email")
			return
		} else if !comparePassword(r.FormValue("password"), user) {
			tmpl.ExecuteTemplate(w, "login.html", "Incorrect Password")
			return
		}
		// Session and Cookie created
		token, expiry := createSession(user)
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   token,
			Expires: expiry,
		})
		tmpl.ExecuteTemplate(w, "index.html", user)
		fmt.Println("Home Loaded", user)
	} else {
		tmpl.ExecuteTemplate(w, "login.html", "User don't exist")
	}

}

// Handler Function to logout
func postLogout(w http.ResponseWriter, r *http.Request) {
	clearCache(w, r)

	// Check if Cookie exist
	c, err := r.Cookie("session_token")
	if err != nil {
		fmt.Println("\nPOST Logout", user)
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
		}
		tmpl.ExecuteTemplate(w, "login.html", "Login to see home")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get Session Token from Cookie
	sessionToken := c.Value

	fmt.Println("\nSessions :", sessions)
	fmt.Println("\nDeleted Session :", sessions[sessionToken])

	// Delete Session from Session Map
	delete(sessions, sessionToken)
	fmt.Println("\nSessions After:", sessions)

	// Set Cookie value as Empty string
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})

	fmt.Println("\nUsers:", Users)

	// Redirect to Login
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

	// Routes
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
