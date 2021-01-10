package httphandlers

import(
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/joho/godotenv"
	"log"
	"strings"
)

const (
	redirectPage = `<a href="https://github.com/login/oauth/authorize?client_id=35cfafe67abea4ab8940&amp;redirect_uri=http://localhost:8080/login/github/callback">Moved Permanently</a>.`
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load("../githubapp.env"); err != nil {
		log.Fatal("No .env file found")
	}
}


func TestRootHandler(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    res := httptest.NewRecorder()

    RootHandler(res, req)

    if res.Code != http.StatusOK {
        t.Errorf("got status %d but wanted %d", res.Code, http.StatusOK)
	}
	
	expected := httpLoginPage
	if strings.TrimSpace(res.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
            res.Body.String(), expected)
	}
}

func TestGithubLoginHandler(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/login/github/", nil)
    res := httptest.NewRecorder()

    GithubLoginHandler(res, req)

    if res.Code != http.StatusMovedPermanently {
        t.Errorf("got status %d but wanted %d", res.Code, http.StatusOK)
	}
	
	expected := redirectPage
	if strings.TrimSpace(res.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
            res.Body.String(), expected)
	}
}

func TestGithubCallbackHandler(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/login/github/callback", nil)
    res := httptest.NewRecorder()

    GithubCallbackHandler(res, req)

    if res.Code != http.StatusOK {
        t.Errorf("got status %d but wanted %d", res.Code, http.StatusOK)
	}
	
	if len(strings.TrimSpace(res.Body.String())) == 0 {
		t.Errorf("handler returned empty body user data")
	}
}