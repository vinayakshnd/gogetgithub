package httphandlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/vinayakshnd/gogetgithub/utils"
)

const (
	httpLoginPage = `<a href="/login/github/">LOGIN</a>`
)

// RootHandler handles `/` request
func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, httpLoginPage)
}

// GithubLoginHandler handles login request `/login/github/`
// for Github auth provider
func GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	githubClientID := utils.GetGithubClientID()

	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=http://localhost:8080/login/github/callback", 
	githubClientID)

	http.Redirect(w, r, redirectURL, 301)
}

// GithubCallbackHandler performs retrieval of github user data after
// successful login
func GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	githubAccessToken := utils.GetGithubAccessToken(code)

	githubData := utils.GetGithubData(githubAccessToken)

	LoggedinHandler(w, r, githubData)
}	

// LoggedinHandler handles callback request 
func LoggedinHandler(w http.ResponseWriter, r *http.Request, githubData string) {
	if githubData == "" {
		// Unauthorized users get an unauthorized message
		fmt.Fprintf(w, "UNAUTHORIZED!")
		return
	}

	w.Header().Set("Content-type", "application/json")

	// Prettifying the json
	var prettyJSON bytes.Buffer
	// json.indent is a library utility function to prettify JSON indentation
	parserr := json.Indent(&prettyJSON, []byte(githubData), "", "\t")
	if parserr != nil {
		log.Panic("JSON parse error")
	}

	// Return the prettified JSON as a string
	fmt.Fprintf(w, string(prettyJSON.Bytes()))
}
