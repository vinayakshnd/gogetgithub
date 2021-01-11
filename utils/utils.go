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
func GetGithubData(accessToken string) string {
	req, reqerr := http.NewRequest("GET", "https://api.github.com/user", nil)
	if reqerr != nil {
		log.Panic("API Request creation failed")
	}

	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	respbody, _ := ioutil.ReadAll(resp.Body)

	return string(respbody)
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

	fileContent := []byte(fmt.Sprintf("This is the content of my file. File updated at %s",
		time.Now().Format(time.RFC1123)))

	opts := &github.RepositoryContentFileOptions{
		Message: github.String("This is my commit message"),
		Content: fileContent,
		Branch:  github.String(githubBranch),
		SHA:     github.String(fileSHA),
		Committer: &github.CommitAuthor{
			Name:  github.String(githubUserFullName),
			Email: github.String(githubUserEmailID),
		},
	}

	// Create or update file file
	_, _, err := ghClient.Repositories.CreateFile(context.Background(), githubUser,
		githubRepoName, githubFileName, opts)
	if err != nil {
		return err
	}
	return nil
}
