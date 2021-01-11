package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/github"
)

const (
	clientIDEnvVar     = "CLIENT_ID"
	clientSecretEnvVar = "CLIENT_SECRET"
)

// GetGithubClientID returns Github OAuth app client ID
// from set environment variable CLIENT_ID
func GetGithubClientID() string {

	githubClientID, exists := os.LookupEnv(clientIDEnvVar)
	if !exists {
		log.Fatal("Github Client ID not defined in .env file")
	}

	return githubClientID
}

// GetGithubClientSecret returns Github OAuth app client secret
// from set environment variable CLIENT_SECRET
func GetGithubClientSecret() string {

	githubClientSecret, exists := os.LookupEnv(clientSecretEnvVar)
	if !exists {
		log.Fatal("Github Client Secret not defined in .env file")
	}

	return githubClientSecret
}

// GetGithubAccessToken returns Github access token
func GetGithubAccessToken(code string) string {

	clientID := GetGithubClientID()
	clientSecret := GetGithubClientSecret()

	requestBodyMap := map[string]string{"client_id": clientID, "client_secret": clientSecret, "code": code}
	requestJSON, _ := json.Marshal(requestBodyMap)

	req, reqerr := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(requestJSON))
	if reqerr != nil {
		log.Panic("Request creation failed")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	respbody, _ := ioutil.ReadAll(resp.Body)

	// Represents the response received from Github
	type githubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	var ghresp githubAccessTokenResponse
	json.Unmarshal(respbody, &ghresp)

	return ghresp.AccessToken
}

// GetGithubData returns Github user data of autorized user
func GetGithubData(accessToken string) (string, error) {
	req, reqerr := http.NewRequest("GET", "https://api.github.com/user", nil)
	if reqerr != nil {
		log.Panic("API Request creation failed")
	}

	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		return "", resperr
	}

	respbody, _ := ioutil.ReadAll(resp.Body)

	return string(respbody), nil
}

// CreateNewBranch creates a new branch from main branch
func CreateNewBranch(ghClient *github.Client,
	githubBranch string,
	githubUser string,
	githubRepoName string) error {

	// Get Commit SHA of tip of master branch
	ref, _, err := ghClient.Git.GetRef(context.Background(), githubUser, githubRepoName, "heads/main")
	if err != nil {
		return err
	}

	ref = &github.Reference{
		Ref: github.String(fmt.Sprintf("heads/%s", githubBranch)),
		Object: &github.GitObject{
			Type: github.String("commit"),
			SHA:  ref.Object.SHA,
		},
	}

	// Create a new branch from main branch
	branch, _, err := ghClient.Git.CreateRef(context.Background(), githubUser, githubRepoName, ref)
	if err != nil {
		return err
	}
	log.Printf("Created new branch ref: [%s]\n", branch.GetRef())
	return nil
}

// CreateNewPullRequest creates a new pull request from
// specified `githubBranch` to `main` branch
func CreateNewPullRequest(ghClient *github.Client,
	githubBranch string,
	githubUser string,
	githubRepoName string) (string, error) {

	pr := &github.NewPullRequest{
		Title: github.String(fmt.Sprintf("Pull request from gogetgithub at %s",
			time.Now().Format(time.RFC1123))),
		Head:                github.String(githubBranch),
		Base:                github.String("main"),
		Body:                github.String("This is pull request raised by gogetgithub app."),
		MaintainerCanModify: github.Bool(true),
	}
	createdPR, _, err := ghClient.PullRequests.Create(context.Background(), githubUser, githubRepoName, pr)
	if err != nil {
		return "", err
	}
	log.Printf("Created new PR: [%s]\n", createdPR.GetHTMLURL())
	return *createdPR.HTMLURL, nil
}

// GetFileBlobSHA returns SHA of given file from given branch
func GetFileBlobSHA(ghClient *github.Client,
	githubBranch string,
	githubUser string,
	githubRepoName string,
	githubFileName string) (string, error) {

	repoContentGetOptions := &github.RepositoryContentGetOptions{
		Ref: githubBranch,
	}

	fContent, _, _, err := ghClient.Repositories.GetContents(context.Background(), githubUser,
		githubRepoName, githubFileName, repoContentGetOptions)

	if err != nil {
		return "", err
	}
	return fContent.GetSHA(), nil
}

// UpdateFile updates given file in specified repository
func UpdateFile(ghClient *github.Client,
	githubBranch string,
	fileSHA string,
	githubUserFullName string,
	githubUserEmailID string,
	githubUser string,
	githubRepoName string,
	githubFileName string) error {

	timestamp := time.Now().Format(time.RFC1123)
	fileContent := []byte(fmt.Sprintf("This is the content of my file. File updated at %s", timestamp))

	opts := &github.RepositoryContentFileOptions{
		Message: github.String(fmt.Sprintf("Commit from gogetgithub at %s", timestamp)),
		Content: fileContent,
		Branch:  github.String(githubBranch),
		SHA:     github.String(fileSHA),
		Committer: &github.CommitAuthor{
			Name:  github.String(githubUserFullName),
			Email: github.String(githubUserEmailID),
		},
	}

	// Create or update file file
	contentCreateResp, _, err := ghClient.Repositories.CreateFile(context.Background(), githubUser,
		githubRepoName, githubFileName, opts)
	if err != nil {
		return err
	}
	log.Printf("Created commit: [%s]\n", *contentCreateResp.Commit.Message)
	return nil
}

// DisplaySuccess function displays success message
func DisplaySuccess(w http.ResponseWriter, githubAccessToken string, prURL string) error {

	w.Header().Set("Content-type", "text/html")
	// Prettifying the json
	var prettyJSON bytes.Buffer
	githubData, err := GetGithubData(githubAccessToken)
	if err != nil {
		return err
	}
	// json.indent is a library utility function to prettify JSON indentation
	parserr := json.Indent(&prettyJSON, []byte(githubData), "", "\t")
	if parserr != nil {
		log.Panic("JSON parse error")
	}

	// Return the prettified JSON as a string
	fmt.Fprintf(w, fmt.Sprintf("\n\nSuccessfully created PR available <a href=\"%s\">here</a>!!!\n\n"+
		"With below Github user:\n\n\n", prURL))
	fmt.Fprintf(w, string(prettyJSON.Bytes()))
	return nil
}
