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

	// TODO: Make redirect URL configurable
	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?scope=repo&client_id=%s&redirect_uri=http://70.0.0.137:20080/login/github/callback",
		githubClientID)

	http.Redirect(w, r, redirectURL, 301)
}

// GithubCallbackHandler performs retrieval of github user data after
// successful login
func GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	githubAccessToken := utils.GetGithubAccessToken(code)

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// list all repositories for the authenticated user
	/*
		repos, _, err := client.Repositories.List(ctx, "", nil)
		fmt.Printf("List of Repos: %+v, Error: %v", repos, err)


		fileContent := []byte("This is the content of my file\nand the 2nd line of it")

		// Note: the file needs to be absent from the repository as you are not
		// specifying a SHA reference here.
		opts := &github.RepositoryContentFileOptions{
			Message:   github.String("This is my commit message"),
			Content:   fileContent,
			Branch:    github.String("master"),
			Committer: &github.CommitAuthor{Name: github.String("Vinayak Shinde"), Email: github.String("vinayakshnd@gmail.com")},
		}
		_, _, err = client.Repositories.CreateFile(ctx, "vinayakshnd", "game-checker", "myNewFile.md", opts)
		if err != nil {
			fmt.Println(err)
			return
		}
	*/

	/*
		newPR := &github.NewPullRequest{
			Title:               github.String("My awesome pull request"),
			Head:                github.String("newbranch"),
			Base:                github.String("master"),
			Body:                github.String("This is the description of the PR created with the package `github.com/google/go-github/github`"),
			MaintainerCanModify: github.Bool(true),
		}

		pr, resp, err := client.PullRequests.Create(context.Background(), "vinayakshnd-org", "game-changer", newPR)
		if err != nil {
			fmt.Println(resp)
			fmt.Println(err)
			return
		}

		fmt.Printf("PR created: %s\n", pr.GetHTMLURL())
	*/

	repo := &github.Repository{
		Name: github.String("newrepo"),
	}

	newRepo, response, err := client.Repositories.Create(context.Background(), "", repo)
	if err != nil {
		fmt.Printf("repo: %+v\nresponse: %+v\n error: %+v", newRepo, response, err)
	}
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
