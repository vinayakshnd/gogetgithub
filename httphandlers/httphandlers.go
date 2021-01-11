package httphandlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/vinayakshnd/gogetgithub/utils"
	"golang.org/x/oauth2"
)

const (
	httpLoginPage      = `<a href="/login/github/">LOGIN</a>`
	githubUser         = "vinayakshnd"
	githubRepoName     = "newrepo"
	githubFileName     = "myNewFile.md"
	githubUserFullName = "Vinayak Shinde"
	githubUserEmailID  = "vinayakshnd@gmail.com"
	githubBranch       = "main"
)

// RootHandler handles `/` request
func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, httpLoginPage)
}

// GithubLoginHandler handles login request `/login/github/`
// for Github auth provider
func GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	githubClientID := utils.GetGithubClientID()

	// TODO: Make redirect URL configurable
	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?scope=repo"+
		"&client_id=%s&redirect_uri=http://70.0.0.137:20080/login/github/callback",
		githubClientID)

	http.Redirect(w, r, redirectURL, 301)
}

// GithubCallbackHandler performs retrieval of github user data after
// successful login
func GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	githubAccessToken := utils.GetGithubAccessToken(code)

	LoggedinHandler(w, r, githubAccessToken)
}

// LoggedinHandler handles callback request
func LoggedinHandler(w http.ResponseWriter, r *http.Request, githubAccessToken string) {
	if githubAccessToken == "" {
		// Unauthorized users get an unauthorized message
		fmt.Fprintf(w, "UNAUTHORIZED!")
		return
	}

	githubData := utils.GetGithubData(githubAccessToken)
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)
	oauthClient := oauth2.NewClient(ctx, ts)
	ghClient := github.NewClient(oauthClient)

	// Get file SHA
	fileSHA, err := utils.GetFileBlobSHA(ghClient, githubBranch, githubUser, githubRepoName, githubFileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = utils.UpdateFile(ghClient, githubBranch, fileSHA, githubUserFullName,
		githubUserEmailID, githubUser, githubRepoName, githubFileName)
	if err != nil {
		fmt.Println(err)
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
	fmt.Fprintf(w, fmt.Sprintf("Successfully updated %s file in https://github.com/%s/%s repository!!!\n"+
		"With below Github user:\n\n\n",
		githubFileName, githubUser, githubRepoName))
	fmt.Fprintf(w, string(prettyJSON.Bytes()))
}
