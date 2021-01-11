package httphandlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

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
	githubBranchPrefix = "main"
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

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)
	oauthClient := oauth2.NewClient(ctx, ts)
	ghClient := github.NewClient(oauthClient)

	githubBranch := fmt.Sprintf("%s_%d", githubBranchPrefix, time.Now().Unix())

	// 1. Create new branch
	err := utils.CreateNewBranch(ghClient, githubBranch, githubUser, githubRepoName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 2. Get blob SHA of file to be updated
	fileSHA, err := utils.GetFileBlobSHA(ghClient, githubBranch, githubUser,
		githubRepoName, githubFileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 3. Update the file and make a commit
	err = utils.UpdateFile(ghClient, githubBranch, fileSHA, githubUserFullName,
		githubUserEmailID, githubUser, githubRepoName, githubFileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 4. Create Pull Request
	prURL, err := utils.CreateNewPullRequest(ghClient, githubBranch, githubUser, githubRepoName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 5. Display Success message
	err = utils.DisplaySuccess(w, githubAccessToken, prURL)
	if err != nil {
		fmt.Println(err)
		return
	}
}
