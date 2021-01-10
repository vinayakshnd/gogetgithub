package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/joho/godotenv"
	"github.com/vinayakshnd/gogetgithub/httphandlers"
	"github.com/vinayakshnd/gogetgithub/utils"
)

// init() executes before the main program
func init() {
	// loads values from .env into the system
	if err := godotenv.Load("githubapp.env"); err != nil {
		log.Fatal("No .env file found")
	}
}

func main() {

	// Simply returns a link to the login route
	http.HandleFunc("/", httphandlers.RootHandler)

	// Login route
	http.HandleFunc("/login/github/", httphandlers.GithubLoginHandler)

	// Github callback
	http.HandleFunc("/login/github/callback", httphandlers.GithubCallbackHandler)

	// Route where the authenticated user is redirected to
	http.HandleFunc("/loggedin", func(w http.ResponseWriter, r *http.Request) {
		httphandlers.LoggedinHandler(w, r, "")
	})

	// Listen and serve on port 8080
	// TODO: make port configurable
	fmt.Println("Listening on port 8080")
	log.Panic(
		http.ListenAndServe(":8080", nil),
	)
}




