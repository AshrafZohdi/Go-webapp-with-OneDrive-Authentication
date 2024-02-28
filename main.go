package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/onedrive"
)

// HomeHandler handles the home route
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the home page!")
}

// AuthHandler initiates the OneDrive authentication flow
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Handle authenticated user (e.g., store user details in a session)
	fmt.Fprintf(w, "Authenticated with OneDrive\nUser ID: %s\nName: %s", user.UserID, user.Name)
}

func main() {
	router := mux.NewRouter()

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	// Initialize OneDrive Router
	onedriveProvider := onedrive.New(
		clientID,
		clientSecret,
		"http://localhost:3000/callback",
	)
	goth.UseProviders(onedriveProvider)

	// Routes for authentication
	router.HandleFunc("/auth/{provider}", AuthHandler)
	router.HandleFunc("/auth/{provider}/callback", CallbackHandler)

	// Home Route
	router.HandleFunc("/", HomeHandler)

	// Start Server
	fmt.Println("Starting Server at Port 3000")
	http.ListenAndServe(":3000", router)
}
