package utils

import (
	"os"
	"testing"
)
func TestGetGithubClientID(t *testing.T) {

    clientID := "abc"
    os.Setenv(clientIDEnvVar, clientID)
    actualClientID := GetGithubClientID()
    expectedClinetID := clientID

    if actualClientID != expectedClinetID {
        t.Errorf("got client id: %s but wanted %s",actualClientID, expectedClinetID)
    }
}

func TestGetGithubClientSecret(t *testing.T) {

    clientSecret := "xyz"
    os.Setenv(clientSecretEnvVar, clientSecret)
    actualClientSecret := GetGithubClientSecret()
    expectedClinetSecret := clientSecret

    if actualClientSecret != expectedClinetSecret {
        t.Errorf("got client id: %s but wanted %s",actualClientSecret, expectedClinetSecret)
    }
}